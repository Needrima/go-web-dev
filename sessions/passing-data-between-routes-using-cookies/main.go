package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"html/template"
	"net/http"
	"os"
)

var tpl *template.Template

type user struct {
	Username, Firstname, Othername string
}

//Both database variables below are global scope so they can be called in different functions
var sessiondb = map[string]string{} // maps session id to user id
var userdb = map[string]user{}      // maps user id to user

func init() {
	tpl = template.Must(template.ParseFiles("index.html", "home.html"))
}

func foo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Session") //Get session cookie
	if err != nil {                    //if cookie does not exist create cookie and set it
		sessionID, err := uuid.NewV4()
		if err != nil {
			fmt.Fprintln(os.Stdout, "Error generating sessionID")
		}
		cookie = &http.Cookie{
			Name:     "Session",
			Value:    sessionID.String(),
			HttpOnly: true,
			//Secure : true,
		}
		http.SetCookie(w, cookie)
	}
	//if cookie exist
	var u user //create user

	// get data from form
	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		f := r.FormValue("firstname")
		o := r.FormValue("othername")
		u = user{un, f, o} //create a user from the form data that is to be passed into the index template on line 48
		//this is where the storage happens
		sessiondb[cookie.Value] = un // map the value of the username to the cookie value since they are both unique
		userdb[un] = u               // map the value of the user info to the username since username is unique
	}
	tpl.ExecuteTemplate(w, "index.html", u) //pass the user info on line 45 to the index template
}

func bar(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Session") // get cookie with this name
	if err != nil {
		http.Redirect(w, r, "/foo", http.StatusTemporaryRedirect)
		return
	}
	//search for the uname using cookie value in the sessions db
	uname, ok := sessiondb[cookie.Value]
	if !ok { //if uname does not exist redirect
		http.Redirect(w, r, "/foo", http.StatusSeeOther)
		return
	}
	user := userdb[uname]                     //search for the user info using the unmae in the user db
	tpl.ExecuteTemplate(w, "home.html", user) //pass the user info into the home index
}

func main() {
	http.HandleFunc("/foo", foo)
	http.HandleFunc("/bar", bar)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}
