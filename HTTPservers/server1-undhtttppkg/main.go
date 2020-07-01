package main

import (
	"fmt"
	"net/http"
	"os"
)

type page string

func (p *page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, _ := os.Create("file.txt")
	f.WriteString("All of una PAPA!")
	w.Header().Write(f)
	fmt.Fprintln(w, "Oya talk make i burst your mouth")
}

func main() {
	var p *page
	http.ListenAndServe(":2020", p)
}
