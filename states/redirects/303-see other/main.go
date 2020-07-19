//http.StatusSeeOther always changes method to post regardless of the initial method.
package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("My root request method", r.Method)
}

func say(w http.ResponseWriter, r *http.Request) {
	fmt.Println("My say request method", r.Method)
	tpl := template.Must(template.ParseFiles("index.html"))
	tpl.Execute(w, r.FormValue("input"))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("My redirect request method", r.Method)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	//or
	//w.Heade.Set("Location", "/")
	//w.WriteHeader(http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/red", redirect)
	http.HandleFunc("/say", say)
	http.ListenAndServe(":8908", nil)
}

//the three paths will all have the method of the root path "/" which is get
//after filling the form, the root path method of will be get.
//"/say" remains get because nothing changed about it
//"/red" becomes post becomes it but changes the method of to root to get cos of http.statusseeother
