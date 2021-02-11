package main

import (
	"os"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("temp.txt"))
}

func main() {
	xs := []int{1, 5, 3, 2, 4}

	tpl.ExecuteTemplate(os.Stdout, "temp.txt", xs)
}
