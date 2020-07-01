package main

import (
	"os"
	"strings"
	"text/template"
)

var fm = template.FuncMap{
	"lc": strings.ToLower,
}

func main() {
	var tpl = template.Must(template.New("").Funcs(fm).ParseGlob("template/*.man"))

	tpl.ExecuteTemplate(os.Stdout, "two.man", "TWO")
}
