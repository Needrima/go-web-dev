package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalln("Error listening:", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln("Error connecting:", err)
		}
		io.WriteString(conn, "Hey!, I'm ready to listen\n")
		fmt.Fprintln(conn, "\t\t oya now i don ready to listen.")
		conn.Close()
	}
}
