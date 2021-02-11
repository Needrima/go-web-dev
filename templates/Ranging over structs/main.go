package main

import (
	"os"
	"text/template"
)

type scientist struct {
	Name, Law string
}

func main() {
	Scientists := []scientist{{"Robert Hooke", "Hooke's law"}, {"Isaac Newton", "Newton's law"}, {"Ohm", "Ohm's law"}}

	var tpl = template.Must(template.ParseFiles("temp.txt", "text.txt"))

	tpl.ExecuteTemplate(os.Stdout, "temp.txt", Scientists)
}
