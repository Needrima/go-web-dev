package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"text/template"
)

var tpl *template.Template

type user struct {
	Sname, Fname, Oname, Email, Username string
	Age                                  int
	Password                             []byte
}

var sessiondb = map[string]string{}
var userdb = map[string]user{}

func init() {
	tpl = template.Must(template.ParseFiles("index.html", "signup.html"))
}

func LoggedIn(r *http.Request) bool { //to check if user is already logged in
	cookie, err := r.Cookie("session")
	if err != nil { //no cookie found; meaning user is not logged in
		return false
	}
	//when user is logged in
	pw := sessiondb[cookie.Value]
	_, ok := userdb[pw]
	return ok // user found; meaning user is logged in
}

func signup(w http.ResponseWriter, r *http.Request) {
	if LoggedIn(r) { //User cannot signup when already logged in
		http.Error(w, "You are logged in already", http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//Process form data
	if r.Method == http.MethodPost {
		sn := r.FormValue("sname")
		fn := r.FormValue("fname")
		on := r.FormValue("oname")
		em := r.FormValue("email")
		un := r.FormValue("uname")
		pw := r.FormValue("pword")
		ag, err := strconv.Atoi(r.FormValue("age")) //convert age to int
		if err != nil {                             //Error occurs when text is passed into age field
			http.Error(w, "Age field allows numbers only", http.StatusInternalServerError)
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		} else if ag < 15 { //Age restriction
			http.Error(w, "Age requirement not met", http.StatusInternalServerError)
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		//encrypt password
		bs, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
		if err != nil {
			fmt.Println("Error encrypting password:", err)
			return
		}

		//check if password already exist
		_, ok := userdb[pw]
		if ok {
			http.Error(w, "Password taken by another user", http.StatusInternalServerError)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		//create cookie and set cookie
		sessionID, _ := uuid.NewV4()
		cookie := &http.Cookie{
			Name:     "session",
			Value:    sessionID.String(),
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		//create user
		u := user{sn, fn, on, em, un, ag, bs}

		//print user info
		fmt.Println("User info:", u)

		//store user info
		sessiondb[cookie.Value] = pw
		userdb[pw] = u

		//redirect back to home page after storing user info
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "signup.html", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func main() {
	http.HandleFunc("/home", Home)
	http.HandleFunc("/signup", signup)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}