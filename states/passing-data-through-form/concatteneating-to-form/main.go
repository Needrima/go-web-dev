package main

import (
	"io"
	"net/http"
)

func foo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	comment := r.FormValue("comment")
	io.WriteString(w, `
		<form action="/" method="POST">
        	<textarea name="comment" id="comm" cols="30" rows="10">
            	Comment Here...
        	</textarea>
        	<button type="submit">Comment</button>
    	</form> <br>`+comment)
}

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":8080", nil)
}

//If the request method is changed to GET, the message will show in the URL after sunmission.
