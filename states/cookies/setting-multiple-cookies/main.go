package main

import (
	"fmt"
	"log"
	"net/http"
)

//to read from a cookie, we use func set cookie from the net/http package
// func SetCookie(w http.ResponseWriter, cookie *Cookie)
//type cookie is a struct that implemets the stringer method
func setmultiplecookies(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "Cookie1",
		Value: "one",
	}) //this sets the cookie

	http.SetCookie(w, &http.Cookie{
		Name:  "Cookie2",
		Value: "Two",
	}) //this sets the cookie

	http.SetCookie(w, &http.Cookie{
		Name:  "Cookie3",
		Value: "three",
	}) //this sets the cookie

	fmt.Fprintln(w, "Cookies has been set. Check dev tools/applications/cookies to view cookie")
}

//to read cookies, we user r.Cookie(cookie name)
func readmultiplecookies(w http.ResponseWriter, r *http.Request) {
	cookie1, err := r.Cookie("Cookie1")
	if err != nil {
		http.Error(w, "No Cookies available for this request", http.StatusBadRequest)
	}
	fmt.Println("cookie:", cookie1)
	fmt.Println("")

	fmt.Fprintln(w, cookie1)

	cookie2, err := r.Cookie("Cookie2")
	if err != nil {
		http.Error(w, "No Cookies available for this request", http.StatusBadRequest)
	}
	fmt.Println("cookie:", cookie2)
	fmt.Println("")

	fmt.Fprintln(w, cookie2)

	cookie3, err := r.Cookie("Cookie3")
	if err != nil {
		http.Error(w, "No Cookies available for this request", http.StatusBadRequest)
	}
	fmt.Println("cookie:", cookie3)

	fmt.Fprintln(w, cookie3)
}

func main() {
	http.HandleFunc("/", setmultiplecookies)
	http.HandleFunc("/read", readmultiplecookies)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	server := http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	log.Fatalln(server.ListenAndServe())
	//the above is same as log.Fatalln(http.ListenAndServe(":8080", nil))
}
