package main

import (
	"html/template"
	"net/http"
	"log"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	err := template.Must(template.ParseFiles("index.html")).Execute(w, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func bar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("AJAX routing with Go"))
}

func main() {
	http.HandleFunc("/foo", bar)
	http.HandleFunc("/", Handle)
	http.ListenAndServe(":8080", nil)
}
