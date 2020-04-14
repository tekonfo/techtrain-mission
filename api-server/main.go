package main

import (
	"database/sql"
	"encoding/json"
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

// これKeyを小文字にすると反応しなくなる、なんで
type CreateUser struct {
	Token string `json:"token"`
}

type GetUser struct {
	Name string `json:"name"`
}

func main() {
	db, err := sql.Open("mysql", "docker:docker@tcp(db:3306)/game_db")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != "GET" {
			fmt.Fprintf(w, "Sorry, only Get methods are supported")
			return
		}

		xtoken, ok := r.URL.Query()["x-token"]
		if !ok || len(xtoken[0]) < 1 {
			log.Println("Url param 'key' is  missing")
		}

		token := xtoken[0]

		// 検索
		row := db.QueryRow("SELECT name FROM user where token = ?;", token)

		if err != nil {
			log.Fatal(err)
		}

		var name string
		err := row.Scan(&name)

		if err != nil {
			fmt.Fprintf(w, "このtokenを持つユーザーは存在しません\n")
			return
		}

		createuser := GetUser{
			Name: name,
		}

		b, err := json.Marshal(createuser)
		if err != nil {
			fmt.Println("error:", err)
		}
		
		w.Write(b)
	})

	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != "POST" {
			fmt.Fprintf(w, "Sorry, only Post methods are supported.")
		}
		
		name := r.FormValue("name")
		t := randSeq(32)
		// ToDo すでに登録されている名前かどうかを判断する
		rows, err := db.Query("SELECT name FROM user WHERE name = ?;", name)

		if err != nil {
			log.Fatal(err)
		}

		if rows.Next() {
			fmt.Fprintf(w, "この名前はすでに登録されています\n")	
		}else {
			_, err := db.Exec("INSERT user (name, token, gacha_times) values (?, ?, 0)", name, t)
			if err != nil {
				log.Fatal(err)
			}

			createuser := CreateUser{
				Token: t,
			}
			b, err := json.Marshal(createuser)
			if err != nil {
				fmt.Println("error:", err)
			}

			w.Write(b)
		}
	})

	fmt.Print("start server")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
