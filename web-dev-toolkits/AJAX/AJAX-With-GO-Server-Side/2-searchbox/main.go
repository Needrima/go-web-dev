package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func main() {
	tpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	})

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		searchResults := []string{"one", "two", "three", "four", "five", "six"}

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Request body reading error:", err)
			return
		}

		fmt.Println("Request body:", string(bs))

		json.NewEncoder(w).Encode(searchResults)
	})

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	http.ListenAndServe(":8080", nil)
}
