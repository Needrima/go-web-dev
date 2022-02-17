package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go.mongo.org/mongo-driver/bson"
	"go.mongo.org/mongo-driver/mongo"
	"go.mongo.org/mongo-driver/mongo/primitive"
	// "go.mongo.org/mongo-driver/mongo/readpref"
	"go.mongo.org/mongo-driver/mongo/options"
)

type autoComplete struct {
	Data []string `bson:"data, omitempty"`
}

var testCollection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://locahost:27017")

	connection, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Printf("Error connecting to mongoDB server: %w\n", err)
		os.Exit(1)
	}

	defer connection.Disconnect(context.TODO())

	if err := connection.Ping(); err != nil {
		log.Printf("Pinging DB failed: %w\n", err)
		os.Exit(1)
	}

	db := connection.Database("gocode")

	testCollection = db.Collection("test")
}

func main() {
	tpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	})

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		var aut autoComplete

		id := primitive.ObjectIdFromHex("6205295f1eea720fd0d0ef68")

		if err := testCollection.FindOne(
			context.TODO(),
			bson.M{"_id": id},
		).Decode(&aut); err != nil {
			log.Println("ERROR FINDING RESULT:", err)
			return
		}

		fmt.Println(aut.Data)
		searchResults := aut.Data

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Request body reading error:", err)
			return
		}

		fmt.Println("Request body:", string(bs))

		json.NewEncoder(w).Encode(searchResults)
	})

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	http.ListenAndServe(":8080", nil)
}
