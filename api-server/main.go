package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

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
// ToDo これってJson返すときに構造体を宣言する以外に方法あるの？
type CreateUser struct {
	Token string `json:"token"`
}

type GetUser struct {
	Name string `json:"name"`
}

func putAccessLog(path string, method string) {
	log.Printf("path: %s, method: %d", path, method)
}

func main() {
	// DB接続
	// ここってなんでdocker名称指定で動くんだっけ
	db, err := sql.Open("mysql", "docker:docker@tcp(db:3306)/game_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Log設定
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// /user
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		putAccessLog(r.URL.Path, r.Method)

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Sorry, only Get methods are supported")
			return
		}

		w.Header().Set("Content-Type", "application/json")

		xtoken, ok := r.URL.Query()["x-token"]
		if !ok || len(xtoken[0]) < 1 {
			log.Println("Url param 'x-token' is  missing")
			w.WriteHeader(http.StatusForbidden)
		}

		token := xtoken[0]

		// 検索
		row := db.QueryRow("SELECT name FROM user where token = ?;", token)
		var name string
		err := row.Scan(&name)

		if err != nil {
			fmt.Fprintf(w, "このtokenを持つユーザーは存在しません\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		createuser := GetUser{
			Name: name,
		}
		b, err := json.Marshal(createuser)
		if err != nil {
			fmt.Println("error:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	// /user/update
	http.HandleFunc("/user/update", func(w http.ResponseWriter, r *http.Request) {
		putAccessLog(r.URL.Path, r.Method)

		w.Header().Set("Content-Type", "application/json")
		if r.Method != "PUT" {
			fmt.Fprintf(w, "Sorry, only Put methods are supported.")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		xtoken, ok := r.URL.Query()["x-token"]
		if !ok || len(xtoken[0]) < 1 {
			log.Println("Url param 'x-token' is  missing")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		token := xtoken[0]
		row := db.QueryRow("SELECT name FROM user WHERE token = ?;", token)

		var oldName string
		err := row.Scan(&oldName)

		if err != nil {
			fmt.Fprintf(w, "このtokenを持つユーザーは存在しません\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		name := r.FormValue("name")
		_, err = db.Exec("UPDATE user SET name = ? WHERE token = ?;", name, token)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	})

	// /user/create
	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		putAccessLog(r.URL.Path, r.Method)

		w.Header().Set("Content-Type", "application/json")
		if r.Method != "POST" {
			fmt.Fprintf(w, "Sorry, only Post methods are supported.")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name")
		t := randSeq(32)
		// ToDo すでに登録されている名前かどうかを判断する
		rows, err := db.Query("SELECT name FROM user WHERE name = ?;", name)

		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if rows.Next() {
			fmt.Fprintf(w, "この名前はすでに登録されています\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		_, err = db.Exec("INSERT user (name, token, gacha_times) values (?, ?, 0)", name, t)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		createuser := CreateUser{
			Token: t,
		}
		b, err := json.Marshal(createuser)
		if err != nil {
			fmt.Println("error:", err)
		}

		w.Write(b)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Print("start server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
