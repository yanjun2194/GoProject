// gotest2 project gotest2.go
package gotest2

import (
	"fmt"
	"net"
	"strconv"
)

func Hello() string {
	words := []string{"hello", "func", "in", "package", "hello"}
	wl := len(words)

	sentence := ""
	for key, word := range words {
		sentence += word
		if key < wl-1 {
			sentence += " "
		} else {
			sentence += "."
		}
	}
	return sentence
}

func ServerBase() {
	fmt.Println("Starting the server...")
	//create listener
	//listener, err := net.Listen("tcp", "192.168.1.27:50000")
	listener, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	// listen and accept connections from clients:
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			return
		}
		//create a goroutine for each request.
		go doServerStuff(conn)
	}

}

func doServerStuff(conn net.Conn) {
	fmt.Println("new connection:", conn.LocalAddr())
	for {
		buf := make([]byte, 20)
		length, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			conn.Close()
			return
		}
		fmt.Println("Receive data from client len:", strconv.Itoa(length), ",str=", string(buf[:length]))
	}
}
