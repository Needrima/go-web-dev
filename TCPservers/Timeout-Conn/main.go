package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func handle(conn net.Conn) {
	if err := conn.SetDeadline(time.Now().Add(time.Second * 15)); err != nil {
		panic("Connection Timeout")
	}
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		scantext := scanner.Text()
		fmt.Println(scantext)
		fmt.Fprintf(conn, "You talk say: %s\n", scantext)
	}
	defer conn.Close()

	//code will get here because a deadline has been established for the connection
	fmt.Println("Code got Here after 15 secs: connection timeout")
}

func main() {
	ln, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalln("Error listening:", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln("Error connecting:", err)
		}
		go handle(conn)
	}
}
