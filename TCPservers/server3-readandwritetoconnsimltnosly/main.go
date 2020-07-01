package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handle(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		scantext := scanner.Text()
		fmt.Println(scantext)
		fmt.Fprintf(conn, "You talk say: %s\n", scantext)
	}
	defer conn.Close()

	//code will never get here because connections are being created continuously
	fmt.Println("Code shouldn't get Here")
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
