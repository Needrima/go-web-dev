package main

import (
	"os"
	"fmt"
	"log"
	"encoding/csv"
	"strconv"
	"github.com/gocolly/colly"
)

func main() {
	file, err := os.Create("data.csv")
	if err != nil {
		log.Println("Creating file", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	collector := colly.NewCollector(colly.AllowedDomains("internshala.com"))

	collector.OnHTML(".internship_meta", func(hm *colly.HTMLElement){
		writer.Write([]string{
			hm.ChildText("a"),
			hm.ChildText("span"),
		})
	})

	for i := 0; i < 312; i++ {
		collector.Visit("https://internshala.com/internships/page-"+ strconv.Itoa(i))
		fmt.Printf("Scraping page %d completed\n", i)
	}
	
	fmt.Println("Scrapped all pages")
}