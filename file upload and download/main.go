package main

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	//"strings"
	"text/template"

	//"time"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := os.Getenv("atlasURI")
	//shellURI := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)

	ctx := context.Background()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("client" + err.Error())
	}
	defer client.Disconnect(ctx)

	database := client.Database("golang")

	tpl := template.Must(template.ParseGlob("templates/*"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseMultipartForm(2 << 20); err != nil {
				http.Error(w, "file exceeds 2 megabytes", 400)
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

			name := r.FormValue("user")

			if err := uploadFile(f, name+"-"+h.Filename, database); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			tpl.ExecuteTemplate(w, "index.html", "File upload successful")
			return
		}

		tpl.ExecuteTemplate(w, "index.html", nil)
	})

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		fileNames, err := downloadFile(database)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(fileNames)

		http.Redirect(w, r, "/check/admin", 303)
	})

	http.HandleFunc("/check/admin", func(w http.ResponseWriter, r *http.Request) {
		var fileNames []string

		dir, err := ioutil.ReadDir("files")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for _, v := range dir {
			fileNames = append(fileNames, v.Name())
		}

		if r.Method == http.MethodGet {
			tpl.ExecuteTemplate(w, "files.html", fileNames)
		} else if r.Method == http.MethodPost {
			filename := r.FormValue("fileName")

			http.ServeFile(w, r, "files/"+filename)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	http.ListenAndServe(":"+port, nil)
}

func uploadFile(file multipart.File, fileNameInDatabase string, db *mongo.Database) error {
	//readfile data
	data, err := ioutil.ReadAll(file)
	if err != nil {
		errors.New("Error occured from upload:" + err.Error())
	}

	// create bucket
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		errors.New("Error occured creating bucket:" + err.Error())
	}

	// upload to bucket stream
	uploadStream, err := bucket.OpenUploadStream(fileNameInDatabase)
	if err != nil {
		errors.New("Error occured creating uploadstream:" + err.Error())
	}
	defer uploadStream.Close()

	// write to upload stream
	filesize, err := uploadStream.Write(data)
	if err != nil {
		errors.New("Error writing to uploadstream:" + err.Error())
	}

	log.Printf("Storing to db successful, Filesize: %d\n", filesize)

	return nil
}

func downloadFile(db *mongo.Database) ([]string, error) {
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
		return []string{}, errors.New("Error finding from fs files: " + err.Error())
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
			return []string{}, errors.New("Error occured creating bucket: " + err.Error())
		}

		// create buffer to write store files
		buf := &bytes.Buffer{}

		// open download stream
		dStream, err := bucket.DownloadToStreamByName(v, buf)
		if err != nil {
			return []string{}, errors.New("DownloadToStreamByName error: " + err.Error())
		}

		f, err := os.Create("files/" + v)
		if err != nil {
			return []string{}, errors.New("Create File error: " + err.Error())
		}
		defer f.Close()

		f.Write(buf.Bytes())

		log.Printf("Download succesful file %v with size %v\n", v, dStream)
	}

	dir, err := ioutil.ReadDir("files")
	if err != nil {
		return []string{}, errors.New("ReadDir error: " + err.Error())
	}

	var names []string

	for _, v := range dir {
		names = append(names, v.Name())
	}

	return names, nil
}
