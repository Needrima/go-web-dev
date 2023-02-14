package main

import (
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"text/template"
)

// create a user struct
type user struct {
	Username, Firstname, Lastname string
	Age                           int
	Password                      []byte
}

// create db variables
var userdb = map[string]user{}
var sessiondb = map[string]string{}

// create template variable
var tpl *template.Template

// function to check if logged in
func loggedin(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	un := sessiondb[cookie.Value]
	_, ok := userdb[un]
	return ok
}

func init() {
	//create template to store html templates
	tpl = template.Must(template.ParseFiles("index.html", "home.html", "login.html"))
	//encrypt password for storage
	pbs, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), bcrypt.MinCost)
	//create a db with a prticular user info
	userdb["Needrima"] = user{"Needrima", "Amirdeen", "Oyebode", 21, pbs}
}

func login(w http.ResponseWriter, r *http.Request) {
	//check if alreaady logged in; u cant login when u are already logged in.
	if loggedin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//process form data to
	if r.Method == http.MethodPost {
		//get usename and password
		pw := r.FormValue("password")
		un := r.FormValue("username")
		//check if user exists using username
		u, ok := userdb[un]
		if !ok {
			http.Error(w, "Username and/or password mismatch", http.StatusForbidden)
			return
		}
		//check if password matches its encrypted hash in the db
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(pw))
		if err != nil {
			http.Error(w, "Username and/or password mismatch", http.StatusForbidden)
			return
		}
		//create cookie
		sessID, _ := uuid.NewV4()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sessID.String(),
		}
		//set the cookie
		http.SetCookie(w, cookie)
		//assign cookie value to the username to keep session alive
		sessiondb[cookie.Value] = un
		//redirect to homepage
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	//check if logged in. i.e homepage can only be accessed when logged in
	if !loggedin(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//if logged in get cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//get user using cookie info
	un := sessiondb[cookie.Value]
	user := userdb[un]
	//pass user info into template
	tpl.ExecuteTemplate(w, "home.html", user)
}

func logout(w http.ResponseWriter, r *http.Request) {
	//cannot logout if not loggedin
	if !loggedin(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//get cookie
	cookie, _ := r.Cookie("session")
	//delete cookie session from database
	delete(sessiondb, cookie.Value)
	//delete cookie to end session
	cookie.MaxAge = -1
	//reset cookie
	http.SetCookie(w, cookie)
	///redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/home", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServe(":8080", nil)
}
