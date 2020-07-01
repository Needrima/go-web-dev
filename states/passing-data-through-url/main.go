package main

import (
	"io"
	"log"
	"net/http"
)

func Intro(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "My name is "+r.FormValue("name")+". and i am "+r.FormValue("age")+" years old")
}

func main() {
	http.HandleFunc("/intro/", Intro)

	err := http.ListenAndServe(":8400", nil)
	if err != nil {
		log.Fatal("Error occured", err)
		return
	}
}

//visit localhost:8400/intro/?name=Amir&age=35
