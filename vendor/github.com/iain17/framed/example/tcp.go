package main

import (
	"net"
	"github.com/iain17/framed"
	"fmt"
	"strings"
	"time"
	"github.com/c2h5oh/datasize"
)

func init() {
	framed.MAX_SIZE = int64(2 * datasize.MB)//Change the size here. Represented in bytes.
}

func server() {
	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":6666")

	// accept connection on port
	conn, _ := ln.Accept()

	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := framed.Read(conn)
		// output message received
		fmt.Print("Message Received:", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(string(message))
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
	}
}

func client() {
	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:6666")
	for {
		time.Sleep(1 * time.Second)
		// listen for reply
		err := framed.Write(conn, []byte("We can't define consciousness because consciousness does not exist.\nHumans fancy that there's something special about the way we perceive the world,\nand yet we live in loops as tight and as closed as the hosts do, seldom questioning our choices, content,\nfor the most part, to be told what to do next.\n- Dr. Ford\n\n"))
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}
		message, err := framed.Read(conn)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}
		fmt.Printf("Message from server:\n%s\n", string(message))
	}
}

func main() {
	go server()
	client()
}