package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

var db *sql.DB

func main() {
	//"root:demopassword@tcp(127.0.0.1:3306)/usersdb" follows the pattern
	//"name_of_connection:password@tcp("connection_host")/database_name"
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/databasename")
	if err != nil {
		fmt.Println("Could not connect to mysql workbench:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/store", Store)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}
func Store(w http.ResponseWriter, r *http.Request) {
	// to input values into database
	if r.Method == http.MethodPost {
		fn := r.FormValue("sname") //type varchar(20) in database
		ln := r.FormValue("lname") //type varchar(20) in database
		id := r.FormValue("ID")    //type int in database
		//signup = name of table in database
		store, err := db.Prepare(`insert into signup values('` + fn + `', '` + ln + `', ` + id + `);`)
		if err != nil {
			fmt.Println("Error:", err)
		}
		defer store.Close()

		n, _ := store.Exec()
		fmt.Println("Executed command:", n)

		tpl.Execute(w, "Sent values to user database")
		return
	}
	tpl.Execute(w, nil)

}
