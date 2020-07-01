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
	books := []string{"QUR'AN", "BIBLE", "TORAH"}
	var tpl = template.Must(template.New("").Funcs(fm).ParseFiles("books.temp"))

	tpl.ExecuteTemplate(os.Stdout, "books.temp", books)
}
