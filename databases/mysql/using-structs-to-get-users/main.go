package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"text/template"
)

type user struct {
	Firstname, Surname string
	ID                 int
}

var tpl *template.Template
var db *sql.DB

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

func main() {
	fmt.Println("My SQL Database Names")
	//"root:demopassword@tcp(127.0.0.1:3306)/usersdb" follows the pattern
	//"name_of_connection:password@tcp("connection_host")/database_name"
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/usersdb")
	if err != nil {
		fmt.Println("Could not connect to mysql workbench:", err)
	}
	defer db.Close()

	http.HandleFunc("/users", Index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("select Firstname, Surname, ID from signup")
	if err != nil {
		fmt.Println(err)
	}

	var users []user

	for rows.Next() {
		var u user
		err := rows.Scan(&u.Firstname, &u.Surname, &u.ID)
		if err != nil {
			fmt.Println(err)
		}

		users = append(users, u)
	}

	for _, us := range users {
		fmt.Println(us.Firstname, us.Surname, us.ID)
	}

	tpl.Execute(w, users)
}
