package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Img(w http.ResponseWriter, r *http.Request) {
	img, err := os.Open("GOW.jpg")
	if err != nil {
		http.Error(w, "Error serving image", http.StatusNotFound)
	}
	defer img.Close()

	_, err = io.Copy(w, img)
	if err != nil {
		fmt.Println("Error copying file", err)
	}
}

func main() {
	http.HandleFunc("/img", Img)
	http.ListenAndServe(":8080", nil)
}
