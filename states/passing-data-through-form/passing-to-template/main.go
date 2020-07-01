package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

type customer struct {
	Name, CID, Status string
}

func foo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	name := r.FormValue("name")
	id := r.FormValue("ID")
	status := r.FormValue("sub") == "on" //checkboxes gives "on" when checked so status is a bool type
	var payment string
	if status {
		payment = "Paid"
	} else {
		payment = "Not paid"
	}

	tpl.Execute(w, customer{name, id, payment})
}

func main() {
	http.HandleFunc("/foo", foo)
	http.ListenAndServe(":8080", nil)
}

//If the request method is changed to GET, the message will show in the URL after sunmission.
