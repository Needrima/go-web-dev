package main

import (
	"log"
	"net/http"
	"golang.org/gorilla/websockets"
)

var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func webSocket(w http.IesponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error occured", err)
		return
	}
	log.Println("Client succesfully connected")

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error occured", err)
		}
		log.Println("Message from conn:", string(msg))

		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println("Error writing back to connection")
			return
		}
		log.Println("Written msg back to connection")

	}
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/ws", webSocket)
	http.ListenAndServe(":8080", nil)
}
