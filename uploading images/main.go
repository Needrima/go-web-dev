package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

// handle concatenation of filenames to cookies
func JoinFileNameToCookie(w http.ResponseWriter, c *http.Cookie, filename string) *http.Cookie {
	s := c.Value
	if !strings.Contains(s, filename) {
		s += "/" + filename
	}
	c.Value = s
	http.SetCookie(w, c)
	return c
}

// chech if a string is present in a slice of strings (for validity of images file extensions)
func Found(xs []string, s string) bool {
	for _, v := range xs {
		if v == s {
			return true
		}
	}
	return false
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
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
		//get file from form
		file, header, err := r.FormFile("img")
		if err != nil {
			fmt.Println("Error getting file:", err)
		}
		defer file.Close()
		//get extension from file
		ext := strings.Split(header.Filename, ".")[1]
		//convert ext to lower case for uniformity
		ext = strings.ToLower(ext)
		//create slice to store image file extensions
		exts := []string{"jpg", "jpeg", "png", "gif"}
		//validate image extension (see func found on line 31)
		if !Found(exts, ext) {
			http.Error(w, "File not an image file", http.StatusBadRequest)
			return
		}
		//create hash to hash file
		//see http://godoc.org/crypto/sha1 for examples on hashes
		hs := sha1.New()
		//copy file into hash
		io.Copy(hs, file)
		//get hashed filename
		hsfilename := fmt.Sprintf("%x", hs.Sum(nil)) + "." + ext
		fmt.Println(hsfilename)
		//get working directory
		WorkingDir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		//create path to file uploads
		PathToUploads := filepath.Join(WorkingDir, "uploads", hsfilename) //remember to create an uploads newfile in the current directory
		//create file using path
		newfile, err := os.Create(PathToUploads)
		if err != nil {
			fmt.Println("Error creating file")
		}
		defer newfile.Close()
		//seek file to start from the beginning
		file.Seek(0, 0)
		//copy file into created file
		io.Copy(newfile, file)
		//join filenames to cookie
		c = JoinFileNameToCookie(w, c, hsfilename)
	}
	//split cookie on joined period
	xs := strings.Split(c.Value, "/")
	fmt.Println(xs)
	//trim out first value to get images since first value is the c.value
	images := xs[1:]
	fmt.Println(images)

	tpl.Execute(w, images)
}

func main() {
	http.HandleFunc("/upload", UploadImage)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads")))) //handle upload folders
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}
