package main

import (
	"net/http"
	"fmt"
	"sync"
	"strings"
	"github.com/gin-gonic/gin"
	"time"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

type request struct {
	Letters string
	Limit int
	Stop bool
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
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

		if message.Stop{
			fmt.Println("Trying to stop")
			stop <- true
			wg.Wait()
		}else{
			message.Letters = strings.ToLower(message.Letters)
			fmt.Println(message.Letters)
			fmt.Println(message.Limit)

			go wordGenerator(conn, message, stop, wg)
			wg.Add(1)
		}
	}
}

func index(c *gin.Context){
	//names := []string{"Robin", "Perez", "Robie", "Robyn", "Cheptileh", "Adhiambo", "Pretty"}
	//name := names[rand.Intn(len(names))]
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"name" : "Robin Cheptileh",
			"year" : time.Now().Year(),
		},
	)
}
