package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func main() {
	//Get tcp address using
	//func ResolveTCPAddr(network, address string) (*TCPAddr, error)
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println("Error getting tcp address:", err.Error())
		return
	}

	//create tcp connection using address
	//func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Error creating connection:", err.Error())
		return
	}
	defer conn.Close()

	///write to connection request some header
	conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))

	//read from connection
	bs, _ := ioutil.ReadAll(conn)
	fmt.Println(string(bs))
}
