// TcpClient project main.go
package main

import (
	"fmt"
	//	"net"
	//	"os"
	//	"strconv"
	//"time"

	"flag"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

//func main() {
//	fmt.Println("start connect!")
//	testClient()
//}

const T_FORMAT = "2006-01-02 15:04:05.000"

// var addr = flag.String("addr", "localhost:12345", "http service address")

var addr = flag.String("addr", "123.206.83.166:2001", "http service address")

func main() {

	for i := 0; i < 10; i++ {
		go startConnect(i)
	}
	select {}
}

func startConnect(index int) {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go timeWriter(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Millisecond * 4000)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format(T_FORMAT)))
	}
}

//func testClient() {
//	con, err := net.Dial("tcp", "39.107.246.237:12345")
//	if err != nil {
//		fmt.Println("connect error")
//		fmt.Println(err.Error())
//		os.Exit(0)
//	}
//	for i := 1; i <= 9; i++ {
//		s := "this is " + strconv.Itoa(i) + "time"
//		// time.Sleep(1000)
//		num, writeerr := con.Write([]byte(s))
//		if writeerr != nil {
//			fmt.Println(writeerr.Error())
//			break
//		}
//		fmt.Println("write length = ", num)
//	}
//	con.Close()
//}
