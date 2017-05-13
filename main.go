package main

import (
	"net/http"
	"time"
	"fmt"
	"strings"
	"sync"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var router *gin.Engine

var wsupgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

type request struct {
	Letters string
	Limit int
	Stop bool
}

type response struct {
	Word string
}

type done struct {
	Done bool
}

func wordGenerator(conn *websocket.Conn, message request, stop chan bool, wg *sync.WaitGroup)  {
	defer wg.Done()
	defer track(time.Now(), "wordGenerator()")

	var temp_word string
	var temp_list []string

	fmt.Println(getPermCount(message.Letters, message.Limit))
	for i := 0; i < getPermCount(message.Letters, message.Limit); i++ {

		select {
		default:
			temp_word = getPerm(message.Letters, message.Limit, i)
			if isInDictionary(temp_word){
				found := false
				for _, v := range temp_list{
					if v == temp_word{
						found = true
					}
				}
				if !found{
					temp_list = append(temp_list, temp_word)
					fmt.Println(temp_word)
					word := response{temp_word}
					err := conn.WriteJSON(word)
					if err != nil {
						return
					}
				}
			}

		case <- stop:
			fmt.Println("Stopped")
			don := done{true}
			err :=  conn.WriteJSON(don)
			if err != nil {
				return
			}
			return
		}
	}

	fmt.Println("Done!")
	don := done{true}
	err :=  conn.WriteJSON(don)
	if err != nil {
		return
	}
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

func getPermCount(letters string, count int) (int) {
	result := 1
	//k characters from a set of n has n!/(n-k)! possible combinations
	for i := len(letters) - count + 1; i <= len(letters); i++ {
		result *= i;
	}
	return result
}

func getPerm(letters string, count int, index int) (string){
	result := ""
	//Decodes index to a $count-length string from $letters, no repeat chars.
	i := 0
	for i < count {
		pos := index % len(letters)
		result += string(letters[int(pos)])
		index = (index - pos) / len(letters)
		a := int(pos + 1)
		b := int(len(letters))
		letters = letters[0:int(pos)] + letters[a:b]
		i = i + 1
	}
	return result
}

func track(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Println()
	fmt.Printf("[TRACKER] : %s took %s", name, elapsed.String())
	fmt.Println()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	gin.SetMode(gin.ReleaseMode)
	// Set the router as the default one provided by Gin
	router = gin.Default()

	router.Static("/css", "templates/css")
	router.Static("/fonts", "templates/fonts")
	router.Static("/js", "templates/js")
	router.Static("/img", "templates/img")

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/index.html")

	// Define the route for the index page and display the index.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	router.GET("/", func(c *gin.Context) {
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

	})

	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	// Start serving the application
	router.Run()

}
