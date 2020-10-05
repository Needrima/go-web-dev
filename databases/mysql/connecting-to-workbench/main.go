package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

var db *sql.DB

func main() {
	//"root:demopassword@tcp(127.0.0.1:3306)/usersdb" follows the pattern
	//"name_of_connection:password@tcp("connection_host")/database_name"
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/usersdb")
	if err != nil {
		fmt.Println("Could not connect to mysql workbench:", err)
	}
	defer db.Close()
	
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/create", CreateTable)
	http.HandleFunc("/delete", DeleteTable)
	http.HandleFunc("/input", Inputdata)
	http.HandleFunc("/get", Getdata)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func CreateTable(w http.ResponseWriter, r *http.Request) { //to create a table in database table
	//func prepare to prepares the sql command
	table, err := db.Prepare("create table users (Surname VARCHAR(40), Firstname VARCHAR(40), ID INT)")
	if err != nil {
		fmt.Println(err)
	}

	defer table.Close()
	//exec() method executes the command
	result, err := table.Exec()
	if err != nil {
		fmt.Println(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintln(w, "Rows affected:", rows)
}

func Inputdata(w http.ResponseWriter, r *http.Request) { // to input values into database
	input, err := db.Query("insert into users values('James', 'Bond', 1)")
	if err != nil {
		fmt.Println(err)
	}

	defer input.Close()
	fmt.Fprintln(w, "Sent values to user database")
}

func Getdata(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("select Firstname from users")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var name string    //var s , name string
	var input []string //s = "Names:\n"

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			fmt.Println("Error:", err)
		}
		input = append(input, name) // s += name +"\n"
	}
	fmt.Fprintln(w, input) //fmt.Fprintln(w, s)
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("drop table users")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	fmt.Fprintln(w, "Drop table users")
}
