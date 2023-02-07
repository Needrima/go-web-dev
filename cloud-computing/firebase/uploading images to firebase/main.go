package main

import (
	"context"
	"fmt"
	"io"
	"os"

	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
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

    objectHandle := bucketHandle.Object("casual.jpg")

    writer := objectHandle.NewWriter(context.Background())
    // very important to set this metadata
    id := uuid.New()
    writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()} 
        defer writer.Close()

    f, err := os.Open("casual.jpg")
    if err != nil {
        fmt.Println("error opening casual.jpg:", err)
        return
    }
    defer f.Close()
    
    if _, err := io.Copy(writer, f); err != nil {
        fmt.Println("error writing to cloud storage wroter:", err)
        return
    }
    
    fmt.Println("all is well")
}