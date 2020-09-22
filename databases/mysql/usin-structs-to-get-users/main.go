package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	Firstname
}

func main() {
	fmt.Println("My SQL Database")
	//"root:demopassword@tcp(127.0.0.1:3306)/usersdb" follows the pattern
	//"name_of_connection:password@tcp("connection_host")/database_name"
	db, err := sql.Open("mysql", "root:Ademola15@tcp(127.0.0.1:3306)/usersdb")
	if err != nil {
		fmt.Println("Could not connect to mysql workbench:", err)
	}

	defer db.Close()

}
