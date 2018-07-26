package main

import (
	"net/http"
	"fmt"
	"sync"
	"strings"
	"github.com/gin-gonic/gin"
	"time"
	"github.com/gorilla/websocket"
	"database/sql"
	_ "github.com/mysql"
	"strconv"
	"os"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type request struct {
	Letters string
	Limit   int
	Stop    bool
}

type apiResponse struct {
	Letters string   `json:"letters"`
	Limit   int      `json:"limit"`
	Words   []string `json:"words"`
	Found   bool     `json:"found"`
}

var DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

func apiHandler(c *gin.Context) {
	//Connect to the database
	db, err := sql.Open("mysql", DSN)
	checkErr(err)
	defer db.Close()

	//Make sure the connection is available
	err = db.Ping()
	checkErr(err)

	stop := make(chan bool)
	found := false
	typ := "api"
	if len(c.Query("letters")) > 0 && len(c.Query("limit")) > 0 {
		let := c.Query("letters")
		lim, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
			return
		}

		req := request{let, lim, false}
		words := apiWordGenerator(req, stop)

		if len(words) > 0 {
			found = true
		}

		apiResp := apiResponse{req.Letters, req.Limit, words, found}
		// insert
		stmt, err := db.Prepare("INSERT requests SET request_type=?,letters=?,letters_limit=?,found=?")
		checkErr(err)
		res, err := stmt.Exec(typ, let, lim, strconv.FormatBool(found))
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		fmt.Println(id)

		c.JSON(http.StatusOK, apiResp)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	checkErr(err)
	fmt.Println("Websocket Initiated")
	defer conn.Close()

	var message request
	stop := make(chan bool)
	wg := &sync.WaitGroup{}

	for {
		err = conn.ReadJSON(&message)
		checkErr(err)
		fmt.Println(message)

		if message.Stop {
			fmt.Println("Trying to stop")
			stop <- true
			wg.Wait()
		} else {
			message.Letters = strings.ToLower(message.Letters)
			fmt.Println(message.Letters)
			fmt.Println(message.Limit)

			go wordGenerator(conn, message, stop, wg)
			wg.Add(1)
		}
	}
}

func index(c *gin.Context) {
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"name": "Robin Cheptileh",
			"year": time.Now().Year(),
		},
	)
}
