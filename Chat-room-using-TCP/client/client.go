package main

import (
	"bufio"
	"fmt"
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

func main() {
	conn, err := net.Dial("tcp", "localhost:8080") //dial conn on ":8080"
	checkError(err, "Erorr dialing connection")    //chech err

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

	conn.Write([]byte(welcomeMsg))
}
