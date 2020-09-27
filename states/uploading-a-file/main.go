package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

func foo(w http.ResponseWriter, r *http.Request) {
	// set header
	w.Header().Set("Content-Type", "text/html; charset=utmpfu-8")

	//declare a varable to convert file to string
	var s string

	if r.Method == http.MethodPost {

		//Call r.FormFile which returns a file, a header and an error
		file, header, err := r.FormFile("upload")
		if err != nil {
			http.Error(w, "Error finding file", http.StatusBadRequest)
		}
		defer file.Close()

		//print variables to terminal
		fmt.Println(file)
		fmt.Println(header)
		fmt.Println(err)

		//call ioutil readall to read from uploaded file
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusNotFound)
		}

		//convert bs to string
		s = string(bs)

		//get file extension
		ext := strings.Split(header.Filename, ".")[1]

		//create tempfile to store file
		tmpfile, err := ioutil.TempFile("uploads", "file-*."+ext)
		if err != nil {
			fmt.Println("Error creating tempfile:", err)
		}
		defer tmpfile.Close()

		tmpfile.Write(bs)
	}
	//write s back to the template
	tpl.Execute(w, s)
}

func main() {
	http.HandleFunc("/foo", foo)
	http.ListenAndServe(":8080", nil)
}
