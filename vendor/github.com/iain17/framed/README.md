# Framed
A super easy to use message framing package for Go. Use this with any object that confirms to io.writer or io.reader. Network connections etc.

## The problem this solves
Please read this [article](https://blog.stephencleary.com/2009/04/message-framing.html) if you don't know what message framing is.
Here is a copy:
> One of the most common beginner mistakes for people designing protocols for TCP/IP is that they assume that message boundaries are preserved. For example, they assume a single “Send” will result in a single “Receive”.

> Some TCP/IP documentation is partially to blame. Many people read about how TCP/IP preserves packets - splitting them up when necessary and re-ordering and re-assembling them on the receiving side. This is perfectly true; however, a single “Send” does not send a single packet.

> Local machine (loopback) testing confirms this misunderstanding, because usually when client and server are on the same machine they communicate quickly enough that single “sends” do in fact correspond to single “receives”. Unfortunately, this is only a coincidence.

> This problem usually manifests itself when attempting to deploy a solution to the Internet (increasing latency between client and server) or when trying to send larger amounts of data (requiring fragmentation). Unfortunately, at this point, the project is usually in its final stages, and sometimes the application protocol has even been published!

This package will make it easy to send data over TCP or a network layer that is reliable. Not UDP, if you want to use UDP please have another layer like UTP on top.


## Limits
I've made the choice to use unsigned 32 bit integer as the message size indication. This means the max size of each write and read will be 2,147,483,647 bytes or 2.15 GB.
This is more than enough for my use cases. If you'd like it even bigger, just fork this repo and change it to uint64.

## Usage
First fetch this package using go get: ```go get github.com/iain17/framed```
Then implement it like so:
```go
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
```