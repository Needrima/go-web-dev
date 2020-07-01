package main

import (
	"log"
	"net/http"
	"net/url"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

type cusreq struct {
	Method           string
	Header           http.Header
	Contentlength    int64
	TransferEncoding []string
	Host             string
	Trailer          http.Header
	RemoteAddr       string
	RequestURI       string
	URL              *url.URL
}

type customhandler bool

func (m customhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqvals := cusreq{
		r.Method,
		r.Header,
		r.ContentLength,
		r.TransferEncoding,
		r.Host,
		r.Trailer,
		r.RemoteAddr,
		r.RequestURI,
		r.URL,
	}
	tpl.ExecuteTemplate(w, "index.html", reqvals)
}

func main() {
	var handle customhandler
	if err := http.ListenAndServe(":8080", handle); err != nil {
		log.Fatalln("Error listening to port:", err)
	}
}
