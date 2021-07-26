package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://e-shop-test:<password>@e-shop-test.csoa6.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("client" + err.Error())
		return
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Ping" + err.Error())
		return
	}

	name, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal("Names" + err.Error())
		return
	}

	for _, v := range name {
		fmt.Println(v)
	}

	fmt.Println("Connected to atlas")

}
