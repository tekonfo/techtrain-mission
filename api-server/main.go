package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"math/rand"

	_ "github.com/go-sql-driver/mysql"
)

// runeってなんやねん
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func main() {
	db, err := sql.Open("mysql", "docker:docker@tcp(db:3306)/game_db")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Fprintf(w, "Sorry, only Post methods are supported.")
		}
		
		name := r.FormValue("name")
		token := randSeq(32)
		// ToDo すでに登録されている名前かどうかを判断する
		rows, err := db.Query("SELECT name FROM user WHERE name = ?;", name)

		if err != nil {
			log.Fatal(err)
		}

		if rows.Next() {
			fmt.Fprintf(w, "この名前はすでに登録されています\n")	
		}else {
			result, err := db.Exec("INSERT user (name, token, gacha_times) values (?, ?, 0)", name, token)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprintf(w, "Name = %s\n", name)
		}
	})

	fmt.Print("start server")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
