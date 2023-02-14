package main

import (
	"context"
	"fmt"
	"io"
	"os"

	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
    "cloud.google.com/go/storage"
)

const (
    bucketName = "your bucket name"
)

func main() {
    opt := option.WithCredentialsFile("serviceAccountKey.json")
    app, err := firebase.NewApp(context.Background(), &firebase.Config{
        StorageBucket: bucketName,
    }, opt)
    
    if err != nil {
        fmt.Printf("error initializing app: %v\n", err)
        return
    }

    client, err := app.Storage(context.TODO())
    if err != nil {
        fmt.Printf("error getting storage client: %v\n", err)
        return
    }

    bucketHandle, err := client.DefaultBucket()
    if err != nil {
        fmt.Printf("error getting storage bucket: %v\n", err)
        return
    }

	objectHandle := bucketHandle.Object("filename")
	attrs, err := objectHandle.Attrs(context.Background())
	if err != nil {
		fmt.Printf("error getting object attributes: %v\n", err)
		return
	}
	
    f, err := os.Create(attrs.Name)
	if err != nil {
		fmt.Printf("error creating file to store image from firebase: %v\n", err)
		return
	}
	defer f.Close()

	reader, err := objectHandle.NewReader(context.Background())
	if err != nil {
		fmt.Printf("error creating reader from storage object: %v\n", err)
		return
	}
	defer reader.Close()

	if _, err := io.Copy(f, reader); err != nil {
		fmt.Println("error writing to cloud storage wroter:", err)
		return
	}

	fmt.Println("image download successful, cheers")
}