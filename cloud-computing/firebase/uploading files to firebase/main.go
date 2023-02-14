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

    f, err := os.Open("casual.jpg")
    if err != nil {
        fmt.Println("error opening casual.jpg:", err)
        return
    }
    defer f.Close()

    objectHandle := bucketHandle.Object(f.Name())

    writer := objectHandle.NewWriter(context.Background())
    // very important to set this metadata for file upload
    id := uuid.New()
    writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()} 

    // set ACL rule to allow all users to have access to file else you can ignore ACLRules
    writer.ObjectAttrs.ACL = []storage.ACLRule{
		{Entity: storage.AllUsers, Role: storage.RoleOwner},
	}
    defer writer.Close()
    
    if _, err := io.Copy(writer, f); err != nil {
        fmt.Println("error writing to cloud storage wroter:", err)
        return
    }
    
    fmt.Println("all is well")
}