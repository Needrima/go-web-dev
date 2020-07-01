package main

import (
	"log"
	"math"
	"os"
	"text/template"
)

var tpl *template.Template

var fm = template.FuncMap{
	"pow": Pow,
	"wn":  math.Round,
	"log": math.Log10,
}

func Pow(x float64) float64 {
	return math.Pow(x, 2)
}

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseFiles("pipe.txt"))
}

func main() {
	if err := tpl.ExecuteTemplate(os.Stdout, "pipe.txt", 4.5); err != nil {
		log.Fatal("error:", err)
	}
}
