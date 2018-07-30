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

func wordGenerator(conn *websocket.Conn, message request, stop chan bool, wg *sync.WaitGroup) {
	// Print DSN URL
	fmt.Println(DSN)

	//Connect to the database
	db, err := sql.Open("mysql", DSN)
	checkErr(err)
	defer db.Close()

	//Make sure the connection is available
	err = db.Ping()
	checkErr(err)
	defer wg.Done()
	defer track(time.Now(), "wordGenerator()")

	var tempWord string
	var tempList []string
	found := false
	typ := "site"

	fmt.Println(getPermutationCount(message.Letters, message.Limit))
	for i := 0; i < getPermutationCount(message.Letters, message.Limit); i++ {

		select {
		default:
			tempWord = getPermutation(message.Letters, message.Limit, i)
			if isInDictionary(tempWord) {
				found := false
				for _, v := range tempList {
					if v == tempWord {
						found = true
					}
				}
				if !found {
					tempList = append(tempList, tempWord)
					fmt.Println(tempWord)
					word := response{tempWord}
					err := conn.WriteJSON(word)
					if err != nil {
						return
					}
				}
			}

		case <-stop:
			fmt.Println("Stopped")
			don := done{true}
			err := conn.WriteJSON(don)
			if err != nil {
				return
			}
			if len(tempList) > 0 {
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
	err = conn.WriteJSON(don)
	if err != nil {
		return
	}
	if len(tempList) > 0 {
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

	var tempWord string
	var tempList []string

	fmt.Println(getPermutationCount(message.Letters, message.Limit))
	for i := 0; i < getPermutationCount(message.Letters, message.Limit); i++ {

		select {
		default:
			tempWord = getPermutation(message.Letters, message.Limit, i)
			if isInDictionary(tempWord) {
				found := false
				for _, v := range tempList {
					if v == tempWord {
						found = true
					}
				}
				if !found {
					tempList = append(tempList, tempWord)
					fmt.Println(tempWord)
				}
			}

		case <-stop:
			fmt.Println("Stopped")
			return tempList
		}
	}

	fmt.Println("Done!")
	return tempList
}

func getPermutationCount(letters string, count int) (int) {
	result := 1
	//k characters from a set of n has n!/(n-k)! possible combinations
	for i := len(letters) - count + 1; i <= len(letters); i++ {
		result *= i
	}
	return result
}

func getPermutation(letters string, count int, index int) (string) {
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
