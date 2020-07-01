package main

import (
	// "html/template"
	"net/http"
	"fmt"
)

// var tpl = template.Must(template.ParseFiles("index.html"))

func Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, `<img src="/inception.jpg">`)
	// tpl.Execute(w, nil)
}

func InceptionImage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "inception.jpg")
}

func main() {
	http.HandleFunc("/", Root) //path must be same with image path on line 13
	http.HandleFunc("/inception.jpg", InceptionImage) //path can also be "/inception"
	//either of the two paths above serves the image
	http.ListenAndServe(":8080", nil)
}
