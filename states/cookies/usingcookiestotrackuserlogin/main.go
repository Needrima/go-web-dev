package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func getcookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("New-cookie") //get cookie on use login
	if err == http.ErrNoCookie {          //ernocookie occurs when no cookie is found
		//whe no cookie s found, set cookie with info below
		cookie = &http.Cookie{
			Name:  "New-cookie",
			Value: "0",
		}
	}

	count, err := strconv.Atoi(cookie.Value) //convert cookie value to int and assign to count to enable increment
	if err != nil {
		log.Fatalln("Error converting cookie.Value", err)
	}

	count++ //increment count i.e for every login or refreshing of page

	cookie.Value = strconv.Itoa(count) //convert count back to string and assign to cookie value

	http.SetCookie(w, cookie) //set cookie back

	fmt.Println("Number of times domain was visited by user: ", count)

	fmt.Fprint(w, "Number of times domain was visited by user: ")
	fmt.Fprintln(w, count)
}

func main() {
	http.HandleFunc("/", getcookie)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	server := http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	log.Fatalln(server.ListenAndServe())
	//the above is same as log.Fatalln(http.ListenAndServe(":8080", nil))
}
