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
		//if error occuered; conn is closed so break loop
		//and send conn to closed channel ie line 47
		if err != nil {
			break
		}
		//if there is a msg
		//range over connetions map
		//item is basically each connection in map for every iteration
		for item := range Conns {
			if conn != item { //every other conn apart from the one writing the message
				item.Write([]byte(msg))
			}
		}
	}
	//loop doesn't run
	//connection is closed
	//send to close conn channel
	closedConn <- conn
}

func main() {
	// Create a  listener
	listener, err := net.Listen("tcp", ":8080")
	checkError(err, "Couldn't create listener")
	defer listener.Close() //close listener before function terminate

	go func() {
		//get conns
		for {
			conn, err := listener.Accept()
			checkError(err, "Error connecting to listener")
			//send connection to newConns channel
			newConns <- conn
			//add conn to connections
			Conns[conn] = true

			defer conn.Close()
		}
	}()

	for {
		select {
		case conn := <-newConns:
			//launching a new gorroutine to increase performance
			//since func WriteToOtherConns is not printing any message
			go WriteToOtherConns(conn)
		case conn := <-closedConn:
			// received on close channel
			for item := range Conns {
				//if close conn is in conns stop the loop
				//and delete the conn from conns
				if conn == item {
					break
				}
			}
			delete(Conns, conn) //inbuilt delete function to delete from maps
		}
	}
}
