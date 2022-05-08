//server side implementation of TCP chat
package main

import (
	"bufio"
	"log"
	"net"
)

var (
	newConns   = make(chan net.Conn) //to store open conns
	closedConn = make(chan net.Conn) //to store closed conns
	Conns      = map[net.Conn]bool{} //or make(map[net.Conn]bool)
	//to store all connections
)

//to check errors
func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(msg+":", err)
	}
}

func WriteToOtherConns(conn net.Conn) {
	//read infinitely from a connection
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		for item := range Conns {
			if conn != item { //every other conn apart from the one writing the message
				item.Write([]byte(msg))
			}
		}
	}

	closedConn <- conn
}

func main() {

	listener, err := net.Listen("tcp", ":8080")
	checkError(err, "Couldn't create listener")
	defer listener.Close() //close listener before function terminate

	go func() {
		//get conns
		for {
			conn, err := listener.Accept()
			checkError(err, "Error connecting to listener")

			newConns <- conn

			Conns[conn] = true

			defer conn.Close()
		}
	}()

	for {
		select {
		case conn := <-newConns:

			go WriteToOtherConns(conn)

		case conn := <-closedConn:

			for item := range Conns {

				if item == conn {
					delete(Conns, item) //delete connection from connections in database
					break
				}

			}
		}
	}
}
