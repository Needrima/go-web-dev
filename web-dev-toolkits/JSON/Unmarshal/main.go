package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type user struct {
	First, Last string
}

func WithUnmarshal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	js := `{"First":"Oyebode","Last":"Amirdeen"}`

	var u user

	err := json.Unmarshal([]byte(js), &u)
	if err != nil {
		log.Println("Error converting json", err)
	}
	
}

func main() {
	http.HandleFunc("/unm", WithUnmarshal)
	http.ListenAndServe(":8080", nil)
}
