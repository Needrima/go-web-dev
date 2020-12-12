package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	Conns       = make(map[*websocket.Conn]bool)
	openConns   = make(chan *websocket.Conn)
	closedConns = make(chan *websocket.Conn)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func Broadcast(conn *websocket.Conn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error occured", err)
			break
		}
		log.Println("Message from conn:", string(msg))

		for item := range Conns {
			if item != conn {
				if err := item.WriteMessage(msgType, msg); err != nil {
					log.Println("Error writing back to connection")
					return
				}
				log.Println("Written msg back to connections")
			}
		}
	}
	closedConns <- conn
}

func webSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	go func() {
		for {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println("Error occured creating connection:", err)
			}
			openConns <- conn

			Conns[conn] = true
		}
	}()

	for {
		select {
		case conn := <-openConns:
			go Broadcast(conn)
		case conn := <-closedConns:
			for item := range Conns {
				if item == conn {
					break
				}
			}
			delete(Conns, conn)
		}
	}
}

func main() {
	http.HandleFunc("/", Home)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/ws", webSocket)
	http.ListenAndServe(":8080", nil)
}
