package main

import (
	"database/sql"
	"os"
	"fmt"
	"log"
	"bufio"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "docker:docker@tcp(db:3306)/game_db")

	// sqlファイルを読み込んで、それを実行する
	fp, err := os.Open("migrate.sql")
	if err != nil {
		log.Fatal(err)
	}

	id := 1
	var name string

	if err := db.QueryRow("SELECT name FROM test_tb WHERE id = ? LIMIT 1", id).Scan(&name); err != nil {
			log.Fatal(err)
	}

	fmt.Println(id, name)
	fmt.Print("test")

	// deferってなんだっけ
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		fmt.Print(scanner.Text())
		// exec
		_, err := db.Exec(scanner.Text()); 
		if err != nil {
			log.Fatal(err)
		}
	}

}