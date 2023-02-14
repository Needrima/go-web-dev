// Package "github.com/satori/go.uuid" is a third party package that has functionalities for generating random values for cookies
package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"text/template"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func bar(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		cookie, err := r.Cookie(r.FormValue("username"))
		if err != nil { //err occurs when there is no cookie available
			uid, err := uuid.NewV4() //uid = unique id randomly generated from package go.uuid
			if err != nil {
				fmt.Printf("Error generating uid: %v\n", err)
			}

			cookie = &http.Cookie{
				Name:  r.FormValue("username"),
				Value: uid.String(),
				//Secure :  true // this line is used when working with https only
				HttpOnly: true, // this line is used when working with http but the cookie cannot be accessed with javascript to increase security.
			}
		}
		http.SetCookie(w, cookie)
		tpl.Execute(w, cookie)
	}
}

func foo(w http.ResponseWriter, r *http.Request) {
	tpl, err := tpl.ParseFiles("delete.html") // adding a new template to var tpl on line 11 with the parsefiles method.
	if err != nil {
		fmt.Println("Error adding template:", err)
	}

	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "delete.html", nil)
	} else if r.Method == http.MethodPost {
		cookie, err := r.Cookie(r.FormValue("delete"))
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/bar", http.StatusTemporaryRedirect)
			return
		}

		cookie.MaxAge = -2 //this line deletes the cookie with the input username
		http.SetCookie(w, cookie)
		tpl.ExecuteTemplate(w, "delete.html", r.FormValue("delete")) //this line gives the info about the cookie deleted
	}
}

func main() {
	http.HandleFunc("/bar", bar)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/foo", foo)
	http.ListenAndServe(":8080", nil)
}
