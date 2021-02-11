package main

import (
	"html/template"
	"net/http"
	"os"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.txt"))
}

func login(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "call.txt", 57)
	tpl.ExecuteTemplate(os.Stdout, "call.txt", nil)
}

func main() {
	http.HandleFunc("/", login)
	http.ListenAndServe(":6060", nil)
}
