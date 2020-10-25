package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type user struct {
	First, Last string
}

func WithMarshal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var u = user{"Oyebode", "Amirdeen"}

	json, err := json.Marshal(u)
	if err != nil {
		log.Println("err")
	}

	w.Write(json)
}

func WithEncoder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var u = user{"Oyebode", "Amirdeen"}
	enc := json.NewEncoder(w)
	err := enc.Encode(u)
	if err != nil {
		log.Println("Error encoding json:", err)
	}
}

func main() {
	http.HandleFunc("/enc", WithEncoder)
	http.HandleFunc("/mar", WithMarshal)
	http.ListenAndServe(":8080", nil)
}
