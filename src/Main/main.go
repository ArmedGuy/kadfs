package main

import (
	"fmt"

	"../d7024e"
)

func main() {
	// TODO
	ip := "127.0.0.1"
	port := 12345
	go d7024e.Listen(ip, port)

	fmt.Scanln()
}
