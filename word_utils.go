package main

import (
	"fmt"
	"sync"
	"time"
	"github.com/gorilla/websocket"
	_ "github.com/mysql"
	"database/sql"
	"strconv"
)

type response struct {
	Word string
}

type done struct {
	Done bool
}

func wordGenerator(conn *websocket.Conn, message request, stop chan bool, wg *sync.WaitGroup)  {
	//Connect to the database
	db, err := sql.Open("mysql", DSN)
	checkErr(err)
	defer db.Close()

	//Make sure the connection is available
	err = db.Ping()
	checkErr(err)
	defer wg.Done()
	defer track(time.Now(), "wordGenerator()")

	var temp_word string
	var temp_list []string
	found := false
	typ := "site"

	fmt.Println(getPermutationCount(message.Letters, message.Limit))
	for i := 0; i < getPermutationCount(message.Letters, message.Limit); i++ {

		select {
		default:
			temp_word = getPermutation(message.Letters, message.Limit, i)
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
			if len(temp_list) > 0 {
				found = true
			}
			stmt, err := db.Prepare("INSERT requests SET request_type=?,letters=?,letters_limit=?,found=?")
			checkErr(err)
			res, err := stmt.Exec(typ, message.Letters, message.Limit, strconv.FormatBool(found))
			checkErr(err)
			id, err := res.LastInsertId()
			checkErr(err)
			fmt.Println(id)
			return
		}
	}

	fmt.Println("Done!")
	don := done{true}
	err =  conn.WriteJSON(don)
	if err != nil {
		return
	}
	if len(temp_list) > 0 {
		found = true
	}
	stmt, err := db.Prepare("INSERT requests SET request_type=?,letters=?,letters_limit=?,found=?")
	checkErr(err)
	res, err := stmt.Exec(typ, message.Letters, message.Limit, strconv.FormatBool(found))
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(id)
}

func apiWordGenerator(message request, stop chan bool) ([]string) {
	defer track(time.Now(), "wordGenerator()")

	var temp_word string
	var temp_list []string

	fmt.Println(getPermutationCount(message.Letters, message.Limit))
	for i := 0; i < getPermutationCount(message.Letters, message.Limit); i++ {

		select {
		default:
			temp_word = getPermutation(message.Letters, message.Limit, i)
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
				}
			}

		case <- stop:
			fmt.Println("Stopped")
			return temp_list
		}
	}

	fmt.Println("Done!")
	return temp_list
}

func getPermutationCount(letters string, count int) (int) {
	result := 1
	//k characters from a set of n has n!/(n-k)! possible combinations
	for i := len(letters) - count + 1; i <= len(letters); i++ {
		result *= i;
	}
	return result
}

func getPermutation(letters string, count int, index int) (string){
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