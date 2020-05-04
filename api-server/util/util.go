package util

import (
	"encoding/json"
	"log"
	"math/rand"
)

// runeってなんやねん
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type ErrorResponse struct {
	// ここのtypeはgrcpのエラーレスポンスを参考にする
	// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
	Type int `json:"int"`
	Message string `json:"message"`
}

func PutAccessLog(path string, method string) {
	log.Printf("path: %s, method: %d", path, method)
}

func GenErrorJson(t int, message string) []byte  {
	error := ErrorResponse{
		Type: t,
		Message: message,
	}
	b, _ := json.Marshal(error)
	return b
}