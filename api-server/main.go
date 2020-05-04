package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"net/http"

	"myself/user"
	"myself/config"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Db接続
	// ここってなんでdocker名称指定で動くんだっけ
	Db, err := sql.Open("mysql", "docker:docker@tcp(db:3306)/game_db")
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	env := &config.Env{Db: Db}

	// Log設定
	f, err := os.OpenFile("logfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	http.HandleFunc("/user", user.GetUserHandler(env))
	http.HandleFunc("/user/update", user.UpdateUserHandler(env))
	http.HandleFunc("/user/create", user.CreateUserHandler(env))

	fmt.Print("start server")

	// ToDo: Graceful shutdownする
	log.Fatal(http.ListenAndServe(":8080", nil))
}
