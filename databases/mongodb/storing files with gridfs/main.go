package main

import (
	"bytes"
	"io/ioutil"
	// "io/ioutil"
	"context"
	"log"

	//"time"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx := context.Background()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("client" + err.Error())
	}
	defer client.Disconnect(ctx)

	database := client.Database("student-devs-blog")

	//uploadFile("./img/myimge.jpg", "newimage.jpg", database)
	downloadFile("newimage.jpg", database)
}

func uploadFile(filePath, fileNameInDatabase string, db *mongo.Database) {
	//readfile data
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Error occured from upload:", err)
	}

	// create bucket
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		log.Fatalln("Error occured creating bucket:", err)
	}

	// upload to bucket stream
	uploadStream, err := bucket.OpenUploadStream(fileNameInDatabase)
	if err != nil {
		log.Fatalln("Error occured creating uploadstream:", err)
	}
	defer uploadStream.Close()

	// write to upload stream
	filesize, err := uploadStream.Write(data)
	if err != nil {
		log.Fatalln("Error writing to uploadstream:", err)
	}

	log.Printf("Storing to db successful, Filesize: %d\n", filesize)
}

func downloadFile(fileNameInDatabase string, db *mongo.Database) {
	// get fs files collection
	fsFiles := db.Collection("fs.files")

	ctx := context.Background()

	// results to show files in fs.files
	results := bson.M{}

	if err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results); err != nil {
		log.Fatalln("Error finding from fs files:", err)
	}
	log.Println(results)

	// get bucket
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		log.Fatalln("Error occured creating bucket:", err)
	}

	// create buffer to write store files
	buf := &bytes.Buffer{}

	// open download stream
	dStream, err := bucket.DownloadToStreamByName(fileNameInDatabase, buf)
	if err != nil {
		log.Fatalln("DownloadToStreamByName error:", err)
	}

	// write to buffer
	if err := ioutil.WriteFile(fileNameInDatabase, buf.Bytes(), 0600); err != nil {
		log.Fatalln("WriteFile error:", err)
	}

	log.Printf("Download succesfull filesie %v\n", dStream)

}
