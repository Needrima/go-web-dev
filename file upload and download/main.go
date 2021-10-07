package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"text/template"

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

	database := client.Database("golang")

	//uploadFile("./files/myimge.jpg", "image.jpg", database)
	//downloadFile(database)

	tpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseMultipartForm(4 << 20); err != nil {
				http.Error(w, "file exceeds 4 megabytes", 400)
				return
			}

			f, h, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "No file chosen", 400)
				return
			}

			ext := filepath.Ext(h.Filename)

			switch ext {
			case ".pdf", ".doc", ".docx":
				log.Println("Accepted file type")
			default:
				http.Error(w, "Unaccepted file type, only .pdf, .doc and .docx files allowed", 400)
				return
			}

			uploadFile(f, h.Filename, database)

			tpl.Execute(w, "File upload successful")
			return
		}

		tpl.Execute(w, nil)
	})

	http.HandleFunc("/admin/check", func(w http.ResponseWriter, r *http.Request) {
		downloadFile(database)
		http.Redirect(w, r, "/", 303)
	})

	http.Handle("/", http.FileServer(http.Dir("files")))

	http.ListenAndServe(":9090", nil)
}

func uploadFile(file multipart.File, fileNameInDatabase string, db *mongo.Database) {
	//readfile data
	data, err := ioutil.ReadAll(file)
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

func downloadFile(db *mongo.Database) {
	type file struct {
		FileName string `bson:"filename"`
	}
	// get fs files collection
	fsFiles := db.Collection("fs.files")

	ctx := context.Background()

	// results to show files in fs.files
	var filenames []string

	cur, err := fsFiles.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalln("Error finding from fs files:", err)
	}

	for cur.Next(ctx) {
		var f file
		cur.Decode(&f)
		filenames = append(filenames, f.FileName)
	}

	for _, v := range filenames {
		// get bucket
		bucket, err := gridfs.NewBucket(db)
		if err != nil {
			log.Fatalln("Error occured creating bucket:", err)
		}

		// create buffer to write store files
		buf := &bytes.Buffer{}

		// open download stream
		dStream, err := bucket.DownloadToStreamByName(v, buf)
		if err != nil {
			log.Fatalln("DownloadToStreamByName error:", err)
		}

		// write to buffer
		f, err := ioutil.TempFile("./files", "*-"+v)
		if err != nil {
			log.Fatalln("TempFile error:", err)
		}
		defer f.Close()
		f.Write(buf.Bytes())

		log.Printf("Download succesful file %v with size %v\n", v, dStream)
	}
}
