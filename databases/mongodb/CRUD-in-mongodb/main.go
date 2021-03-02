package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
)

type User struct { //exported field in capital letters
	First, Last, Gender string
}

var session *mgo.Session
var collection *mgo.Collection

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	var err error
	session, err = mgo.Dial("mongodb://localhost:27017")

	if err != nil {
		log.Println("Couldnt connect", err)
	}
	defer session.Close()

	collection = session.DB("gomongo").C("go-web-dev")

	http.HandleFunc("/create", Create)
	http.HandleFunc("/read", Read)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Remove)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "create.html", nil)
	} else if r.Method == http.MethodPost {
		fn := r.FormValue("fname")
		ln := r.FormValue("lname")
		gn := r.FormValue("sex")

		err := collection.Insert(&User{fn, ln, gn})
		if err != nil {
			fmt.Println("Error creating user")
			return
		}
		fmt.Println("User created successfully")

		tpl.ExecuteTemplate(w, "create.html", "User created successfully")
	}
}

func Read(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "read.html", nil)
	} else if r.Method == http.MethodPost {
		fn := r.FormValue("fname")

		var u User

		err := collection.Find(bson.M{"first": fn}).One(&u)
		if err != nil {
			fmt.Println("Error finding user")
			return
		}
		fmt.Println("User found")

		tpl.ExecuteTemplate(w, "read.html", u)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "update.html", nil)
	} else if r.Method == http.MethodPost {
		fn := r.FormValue("fname")
		ln := r.FormValue("lname")

		err := collection.Update(bson.M{"first": fn}, bson.M{"$set": bson.M{"last": ln}})
		if err != nil {
			fmt.Println("Error finding user")
			return
		}
		fmt.Println("User updated")

		tpl.ExecuteTemplate(w, "update.html", "user updated")
	}
}

func Remove(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "delete.html", nil)
	} else if r.Method == http.MethodPost {
		fn := r.FormValue("fname")

		err := collection.Remove(bson.M{"first": fn})
		if err != nil {
			fmt.Println("Error deleting user")
			return
		}
		fmt.Println("User deleted")

		tpl.ExecuteTemplate(w, "delete.html", "user deleted")
	}
}
