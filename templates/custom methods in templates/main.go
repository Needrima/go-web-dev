package main

import (
	"log"
	"os"
	"text/template"
)

type person struct {
	Name string
	Age  int
}

func (p person) Actualage(n int) int {
	return n * p.Age
}

func (p person) Fullname(s string) string {
	return s + " " + p.Name
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

func main() {
	p1 := person{"Amir", 10}
	if err := tpl.ExecuteTemplate(os.Stdout, "tpl.gohtml", p1); err != nil {
		log.Fatalln("Error parsing template:", err)
	}
}
