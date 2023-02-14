package main

import (
	"bufio"
	"fmt"
	//"github.com/gookit/color"
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

// read messages from other connections and print to terminal
func ReadFromConnection(conn net.Conn) {
	//infinite for to run process forever
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')

		if err == io.EOF {
			conn.Close()
			log.Println("Connection closed")
			os.Exit(1)
		}

		fmt.Println(msg)
		fmt.Println("------------------------------------------------")
	}
}

// reads from terminal and writes to connection
func WriteToConnection(conn net.Conn, name string) {
	for {
		reader := bufio.NewReader(os.Stdin) // read from terminal
		msg, err := reader.ReadString('\n') // get message

		fmt.Println("\n------------------------------------------------")

		if err != nil {
			break
		}

		msg = fmt.Sprintf("%s:- %s\n", name, strings.Trim(msg, " \r\n"))

		conn.Write([]byte(msg))
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080") //dial conn on ":8080"
	checkError(err, "Erorr dialing connection")    //check err

	defer conn.Close() //close conn before function terminate
	//get name

	fmt.Print("Enter a room name: ")

	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	checkError(err, "Could not get name")
	roomName = strings.Trim(name, " \r\n")

	welcomeMsg := fmt.Sprintf("Welcome %s, Chat with friends.\n", name)

	fmt.Println(welcomeMsg)

	go ReadFromConnection(conn)

	WriteToConnection(conn, name)
}
