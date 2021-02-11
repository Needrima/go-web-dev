package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// declare upgrader struct
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // not important
	WriteBufferSize: 1024, // not important
}

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

//websocket function
func webSocket(w http.ResponseWriter, r *http.Request) {
	//chechorigin acepts all incoming connections
	//good practicce to be called
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	//create connection with upgrager by calling upgrade method
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error occured", err)
		return
	}
	defer conn.Close()

	log.Println("Client succesfully connected")

	//read messages continuously from connetion and write message back to connetion at "client.html"
	for {
		//read
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error occured reading msg:", err)
		}
		log.Println("Message from conn:", string(msg))
		//write
		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println("Error writing back to connection")
			return
		}
		log.Println("Written msg back to connection")
	}
}

func main() {
	http.HandleFunc("/", Home)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/ws", webSocket)
	http.ListenAndServe(":8080", nil)
}
