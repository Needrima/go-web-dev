package main

import (
	"fmt"
	"log"
	"net/http"
)

func customservemux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "DON'T LEAVE ME!")
}

func customservemux2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "TAKE ME WITH YOU!")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/cust1/", customservemux)
	mux.HandleFunc("/cust2/", customservemux2)
	//or
	// mux.Handle("/cust1/", http.HandlerFunc(customservemux))
	// mux.Handle("/cust2/", http.HandlerFunc(customservemux2))
	// http.HandlerFunc is used to convert functions to handlers which mux.Handle needs

	log.Fatalln(http.ListenAndServe(":8080", mux))
}
