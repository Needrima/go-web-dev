package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(msg+":", err)
	}
}

//read messages from other connections and print to terminal
func ReadFromConnection(conn net.Conn) {
	//infinite for to run process forever
	for {
		reader := bufio.NewReader(conn)     // read from connection
		msg, err := reader.ReadString('\n') // get message

		if err == io.EOF { //no message on connections
			conn.Close()                     //close conn
			log.Fatalln("Connection closed") //alert connection close and exit program
		}

		fmt.Println(msg)
		fmt.Println("------------------------------------------------")
	}
}

//write messages through teminal to other conections
func WriteToConnection(conn net.Conn, roomName string) {
	for {
		reader := bufio.NewReader(os.Stdin) // read from terminal
		msg, err := reader.ReadString('\n') // get message

		if err != nil {
			break
		}

		msg = fmt.Sprintf("%s:- %s\n", roomName, strings.Trim(msg, " \r\n"))

		conn.Write([]byte(msg))
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080") //dial conn on ":8080"
	checkError(err, "Erorr dialing connection")    //chech err
	defer conn.Close()

	defer conn.Close() //close conn before function terminate
	//get name

	fmt.Print("Enter a room name: ")

	reader := bufio.NewReader(os.Stdin)        //read username from terminal
	roomName, err := reader.ReadString('\n')   // assign to variable
	checkError(err, "Could not get roomname")  //chech err
	roomName = strings.Trim(roomName, " \r\n") // trim out whitespaces from name

	welcomeMsg := fmt.Sprintf("Welcome %s, Chat with friends.\n", roomName) //notice the trailing "\n"
	//absence of it returns an EOF error on the sever side
	//when reading the message from the connectio
	fmt.Println(welcomeMsg)

	go ReadFromConnection(conn)

	WriteToConnection(conn, roomName)
}
