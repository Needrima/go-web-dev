package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

func JoinFileNameToCookie(w http.ResponseWriter, c *http.Cookie, filename string) *http.Cookie {
	s := c.Value
	if !strings.Contains(s, filename) {
		s += "/" + filename
	}
	c.Value = s
	http.SetCookie(w, c)
	return c
}

func Found(xs []string, s string) bool {
	for _, v := range xs {
		if v == s {
			return true
		}
	}
	return false
}

func foo(w http.ResponseWriter, r *http.Request) {
	//create cookie
	c, err := r.Cookie("images")
	if err != nil {
		Sid, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "images",
			Value: Sid.String(),
		}
		http.SetCookie(w, c)
	}

	if r.Method == http.MethodPost {

		//Call r.FormFile which returns a file, a header and an error
		file, header, err := r.FormFile("upload")
		if err != nil {
			http.Error(w, "Error finding file", http.StatusBadRequest)
		}
		defer file.Close()

		//call ioutil readall to read from uploaded file
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusNotFound)
		}

		//split file to get extension
		splitFile := strings.Split(header.Filename, ".")

		//get file extension
		ext := splitFile[len(splitFile)-1]
		ext = strings.ToLower(ext)
		//check if extension is an image file

		exts := []string{"jpg", "jpeg", "png", "gif"}

		if !Found(exts, ext) {
			http.Error(w, "File not an image file", http.StatusInternalServerError)
			return
		}

		//create tempfile to store file
		tmpfile, err := ioutil.TempFile("uploads", "img-*."+ext)
		if err != nil {
			fmt.Println("Error creating tempfile:", err)
		}
		defer tmpfile.Close()

		tmpfile.Write(bs) //or io.Copy(tmpfile, file)

		files, err := ioutil.ReadDir("uploads")
		if err != nil {
			log.Println("Could not read uploads directory:", err)
			return
		}

		for _, file := range files {
			c = JoinFileNameToCookie(w, c, file.Name())
		}
	}
	xs := strings.Split(c.Value, "/")
	fmt.Println(xs)
	//trim out first value to get images since first value is the c.value
	images := xs[1:]
	fmt.Println(images)

	tpl.Execute(w, images)
}

func main() {
	http.HandleFunc("/foo", foo)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	http.ListenAndServe(":8080", nil)
}
