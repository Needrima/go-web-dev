package main

import (
	"log"
	"net/http"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

type mux int

func (m mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tpl.ExecuteTemplate(w, "index.html", r.Form)
}

func main() {
	var p mux
	if err := http.ListenAndServe(":8080", p); err != nil {
		log.Fatalln(err)
	}
}
