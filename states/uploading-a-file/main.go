package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

func foo(w http.ResponseWriter, r *http.Request) {
	// set header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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

		// to write uploaded files back to a file
		filecontent := filepath.Join("./useruploads/", header.Filename)

		//create a new file to store filecontent
		filedst, err := os.Create(filecontent)
		if err != nil {
			http.Error(w, "Error Creating file to store uploads", http.StatusNotFound)
		}
		defer filedst.Close()
		filedst.Write(bs)
	}
	//write s back to the template
	tpl.Execute(w, s)
}

func main() {
	http.HandleFunc("/foo", foo)
	http.ListenAndServe(":8080", nil)
}
