package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func ShoutAyo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Shout Ayo!!!")
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}
func ShoutName(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET"{
		tpl.ExecuteTemplate(w, "index.html", nil)
	}else {
		tpl.ExecuteTemplate(w, "index.html", r.FormValue("name"))
	}
}

func main() {
	// Most elegant procedure for routing
	http.HandleFunc("/", ShoutAyo)
	http.HandleFunc("/name", ShoutName)
	// or
	// http.Handle("/cust1/", http.HandlerFunc(customservemux))
	// http.Handle("/cust2/", http.HandlerFunc(customservemux2))
	// http.HandlerFunc is used to convert functions to handlers which mux.Handle needs


	log.Fatalln(http.ListenAndServe(":8080", nil))
}
