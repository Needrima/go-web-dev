package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func Image(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<img src="/God of war.jpg">`)
}

func Gearimage(w http.ResponseWriter, r *http.Request) {
	opengear, _ := os.Open("Gear of war.jpg")
	defer opengear.Close()
	io.Copy(w, opengear)
}

func Gowimage(w http.ResponseWriter, r *http.Request) {
	opengow, _ := os.Open("God of war.jpg")
	defer opengow.Close()
	io.Copy(w, opengow)
}

func main() {
	http.HandleFunc("/", Image)             // serves God of war image
	http.HandleFunc("/God of war.jpg", Gowimage) // serves God of war image
	http.HandleFunc("/gearimg/", Gearimage) //when called, serves Gear of war image
	
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
