package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type autoComplete struct {
	Data []string `json:"data", bson:"data"`
}

var testCollection *mgo.Collection

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Session:", err)
	}

	if err := session.Ping(); err != nil {
		log.Printf("Pinging DB failed: %w\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to mongoDB")

	db := session.DB("gocode")

	testCollection = db.C("test")

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	router.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {
		hex := "6205295f1eea720fd0d0ef68"
		id := bson.ObjectIdHex(hex)

		fmt.Println("ObjectId:", id)

		var aut autoComplete

		if err := testCollection.Find(bson.M{"_id": id}).One(&aut); err != nil {
			log.Println("Error finding document:", err)
			return
		}

		fmt.Printf("Document: %v, Data: %v\n", aut, aut.Data)

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			return
		}

		fmt.Println("Request body:", string(bs))

		json.NewEncoder(w).Encode(aut)
	})

	fs := http.FileServer(http.Dir("./public"))
	router.Handle("/public/", http.StripPrefix("/public/", fs))

	log.Fatal(http.ListenAndServe(":8080", router))
}
