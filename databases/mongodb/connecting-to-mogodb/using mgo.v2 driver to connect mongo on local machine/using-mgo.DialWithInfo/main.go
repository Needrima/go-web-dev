package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type result struct {
	Name   string
	Matric int
	Score  float64
}

func main() {
	dialinfo, err := mgo.ParseURL("mongodb://localhost:27017")
	if err != nil {
		log.Println("Dialinfo:", err)
		return
	}

	fmt.Println(dialinfo)

	session, err := mgo.DialWithInfo(dialinfo)
	if err != nil {
		log.Println("Dialinfo:", err)
		return
	}
	defer session.Close()

	collection := session.DB("golang").C("go_practice")

	err = collection.Insert(&result{"Fawaz", 160809034, 28})
	if err != nil {
		log.Println("Inserting:", err)
		return
	}

	var r []result
	err = collection.Find(bson.M{}).All(&r)
	if err != nil {
		log.Println("Find:", err)
		return
	}

	for _, v := range r {
		fmt.Println("Score:", v.Score)
	}
}
