package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	//Get tcp address using
	//func ResolveTCPAddr(network, address string) (*TCPAddr, error)
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println("Error getting tcp address:", err.Error())
		return
	}

	//create a listener
	//func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Error creating listener:", err.Error())
		return
	}

	//create conn from listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting to listener", err)
			continue
		}
		defer conn.Close()
		//write to connection
		time := time.Now().String()
		io.WriteString(conn, time) //or conn.Write([]byte(time))
	}
}

//run run program and visit localhost:8080
