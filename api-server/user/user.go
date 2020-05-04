package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"myself/util"
	"myself/config"
	)

// Q これってJson返すときに構造体を宣言する以外に方法あるの？
// A 基本的にはない、名前なしで簡単に構造体宣言する方法もあるが、だいたい使わない
// Res,Req両方の構造体を宣言する方法が一般的らしい、なるほどな
type CreateUser struct {
	Token string `json:"token"`
}

type GetUser struct {
	Name string `json:"name"`
}

func CreateUserHandler(env *config.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.PutAccessLog(r.URL.Path, r.Method)

		w.Header().Set("Content-Type", "application/json")
		if r.Method != "POST" {
			fmt.Fprintf(w, "Sorry, only Post methods are supported.")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name")
		t := util.RandSeq(32)
		// ToDo すでに登録されている名前かどうかを判断する
		rows, err := env.Db.Query("SELECT name FROM user WHERE name = ?;", name)

		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if rows.Next() {
			fmt.Fprintf(w, "この名前はすでに登録されています\n")
			// 401の方がよく使われるらしい
			// 世のエラーはだいたいエラータイプを自分で作成するらしい
			w.WriteHeader(http.StatusForbidden)
			return
		}

		_, err = env.Db.Exec("INSERT user (name, token, gacha_times) values (?, ?, 0)", name, t)
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
}

func GetUserHandler(env *config.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.PutAccessLog(r.URL.Path, r.Method)
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			errorJson := util.GenErrorJson(12, "Method not found at server")
			w.Write(errorJson)
			return
		}

		xtoken, ok := r.URL.Query()["x-token"]
		if !ok || len(xtoken[0]) < 1 {
			w.WriteHeader(http.StatusForbidden)
			errorJson := util.GenErrorJson(3, "Url param 'x-token' is  missing")
			w.Write(errorJson)
			return
		}

		token := xtoken[0]
		// SQLインジェクション対策をする
		// escape処理をする
		row := env.Db.QueryRow("SELECT name FROM user where token = ?;", token)
		var name string
		err := row.Scan(&name)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			errorJson := util.GenErrorJson(3, "このトークンを持つユーザーは存在しません。")
			w.Write(errorJson)
			return
		}

		createuser := GetUser{
			Name: name,
		}

		b, err := json.Marshal(createuser)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			errorJson := util.GenErrorJson(3, "Url param 'x-token' is  missing")
			w.Write(errorJson)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
}

func UpdateUserHandler(env *config.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.PutAccessLog(r.URL.Path, r.Method)
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
		row := env.Db.QueryRow("SELECT name FROM user WHERE token = ?;", token)

		var oldName string
		err := row.Scan(&oldName)

		if err != nil {
			fmt.Fprintf(w, "このtokenを持つユーザーは存在しません\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		name := r.FormValue("name")
		_, err = env.Db.Exec("UPDATE user SET name = ? WHERE token = ?;", name, token)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}