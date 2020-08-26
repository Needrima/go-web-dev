package main

import (
	"fmt"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("My root request method", r.Method)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("My redirect request method", r.Method)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	//or
	//w.Heade.Set("Location", "/")
	//w.WriteHeader(http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/red", redirect)
	http.ListenAndServe(":8908", nil)
}

//running /red gives the req method for both paths first
//running a second time gives the req method for "/"" only cos of statusmovedpermanently
