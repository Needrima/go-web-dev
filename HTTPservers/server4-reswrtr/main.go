package main

import (
	"log"
	"net/http"
)

type customhandler bool

func (m customhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
}

func main() {
	var handle customhandler
	if err := http.ListenAndServe(":8080", handle); err != nil {
		log.Fatalln("Error listening to port:", err)
	}
}
