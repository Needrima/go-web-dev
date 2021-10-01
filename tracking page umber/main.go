package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var pageNumber = 0

var tpl *template.Template

func inc(x int) int {
	x++
	return x
}

func dec(x int) int {
	x--
	return x
}

var fm = template.FuncMap{
	"next": inc,
	"prev": dec,
}

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseFiles("index.html"))
}

func main() {
	
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "index.html", pageNumber)
	})

	http.HandleFunc("/next/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/next/"):]
		pageNumber, _ := strconv.Atoi(path)
		tpl.ExecuteTemplate(w, "index.html", pageNumber)
	})

	http.HandleFunc("/previous/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/previous/"):]
		pageNumber, _ := strconv.Atoi(path)
		if pageNumber < 0 {
			http.Error(w, "end of pages", 400)
			return
		}

		tpl.ExecuteTemplate(w, "index.html", pageNumber)
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}
