package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			scanner := bufio.NewScanner(conn)
			scanner.Scan()
			fmt.Println(scanner.Text())
			conn.Write(bytes.ToUpper(scanner.Bytes()))
			conn.Close()
		}()
	}

}
