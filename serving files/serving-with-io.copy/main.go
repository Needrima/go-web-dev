package main

import (
	"io"
	"log"
	"net/http"
)

func Image(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<img src="/God of war.jpg">`)
}

func Gowimage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "God of war.jpg")
}

func main() {
	http.HandleFunc("/", Image) //when called, serves God of war image
	http.HandleFunc("/God of war.jpg", Gowimage)
	http.HandleFunc("/Gow", Gowimage) //when called, also serves God of war image
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
