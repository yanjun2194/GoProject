// gotest project main.go
package main

import (
	"fmt"
	"gotest2"
)

func main() {
	var i int = 20
	i += 30
	fmt.Println("Hello Test!")
	bb := gotest2.Hello()
	//gotest2.ServerBase()

	gotest2.StartListening()

	fmt.Println(bb)
	fmt.Println("Hello World!")
}
