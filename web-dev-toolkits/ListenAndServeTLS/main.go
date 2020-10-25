package main

import (
	"fmt"
	"net/http"
)

func foo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "Hello from TLS..... This is a secure connection")
}

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
}

//to generate cert.pem and key.pem
//visit C:\Go\src\crypto\tls
//copy generate_cert.go and run the program specifying host
//go run generate_cert.go --host=localhost or --host=somedomainame.com
//copy cert.pem and key.pem to workspace
//Make sure to delete key.pem and cert.pem before pushing to github or any other source control platform
//go run main.go and visit https://localhost:10443
//chrome might give a warning...Just by-pass 