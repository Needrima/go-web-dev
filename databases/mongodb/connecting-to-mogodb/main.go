package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type User struct {
	First, Last, Gender string
}

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")

	if err != nil {
		log.Println("Couldnt connect", err)
	}
	defer session.Close()

	collection := session.DB("users").C("customers")

	err = collection.Insert(&User{"Fawaz", "oyebode", "male"})
	if err != nil {
		log.Println("Couldnt insert data", err)
	}

	var u User

	err = collection.Find(bson.M{"first": "Fawaz"}).One(&u)
	if err != nil {
		log.Println("Couldnt find data", err)
	}

	fmt.Println(u)
}
