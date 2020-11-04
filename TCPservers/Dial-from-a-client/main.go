package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8083")
	if err != nil {
		fmt.Println("Error dialing client", err.Error())
	}
	defer conn.Close()

	//fmt.Fprintln(conn, "Client dialed")

	fmt.Println("Enter a text: ")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			fmt.Println("Empty text field")
		} else {
			fmt.Fprintln(conn, text)
			//conn.Write([]byte(text))
		}

		fmt.Print("Enter a text: ")
	}
}

//open two terminals
//run the code here on one
//run the code in Read-from-TCP-connection in the other
