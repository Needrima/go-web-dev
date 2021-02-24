package main

import (
	"fmt"
	"net/http"
)

func foo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I love %s\n", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":8080", nil)
}
