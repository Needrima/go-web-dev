package main

import (
	"fmt"
	"log"
	"net/http"
)

type customservemux int

func (c customservemux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "From customservemux one")
}

type customservemux2 int

func (c2 customservemux2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "From customservemux two")
}

func main() {
	var c customservemux
	var c2 customservemux2

	mux := http.NewServeMux()
	mux.Handle("/cust1/", c)
	mux.Handle("/cust2/", c2)

	log.Fatalln(http.ListenAndServe(":8080", mux))
}
