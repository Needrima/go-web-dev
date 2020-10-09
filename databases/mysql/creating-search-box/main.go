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
	var err error
	db, err = sql.Open("mysql", "root:Ademola15@tcp(127.0.0.1:3306)/usersdb")
	if err != nil {
		fmt.Println("Could not connect to mysql workbench:", err)
	}
	defer db.Close()

	err = db.Ping() 
	if err != nil {
		fmt.Println("error, pinging", err)
	}

	http.HandleFunc("/search", Search)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.Execute(w, nil)
	}else {
	// name := r.FormValue("search")

	query := `select * from signup where Firstname = "Amirdeen";`

	row := db.QueryRow(query)

	var us user

	err := row.Scan(&us.Firstname, &us.Surname, &us.ID)
	if err != nil {
		fmt.Println("Error getting row", err)
	}

	tpl.Execute(w, us)
	}	
}
