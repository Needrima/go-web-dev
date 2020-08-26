package main

import (
	"fmt"
	"log"
	"net/http"
)

//to read from a cookie, we use func set cookie from the net/http package
// func SetCookie(w http.ResponseWriter, cookie *Cookie)
//type cookie is a struct that implemets the stringer method
func set(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "New-cookie",
		Value: "My-cookie",
	}) //this sets the cookie

	fmt.Fprintln(w, "Cookie has been set. Check dev tools/applications/cookies to view cookie")
}

//to read cookies, we user r.Cookie(cookie name)
func read(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("New-cookie")
	if err != nil {
		http.Error(w, "No Cookies available for this request", http.StatusBadRequest)
	}
	fmt.Println("cookie:", cookie)

	fmt.Fprintln(w, cookie)
}

func main() {
	http.HandleFunc("/", set)
	http.HandleFunc("/read", read)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	server := http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	log.Fatalln(server.ListenAndServe())
	//the above is same as log.Fatalln(http.ListenAndServe(":8080", nil))
}
