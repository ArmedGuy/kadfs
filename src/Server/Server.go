package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"../Messages"
	"github.com/golang/protobuf/proto"
)

func Listen(ip string, port int) {
	s := ip + ":" + strconv.Itoa(port)

	addr, err := net.ResolveUDPAddr("udp", s)
	checkError(err)

	conn, err := net.ListenUDP("udp", addr)
	checkError(err)

	buffer := make([]byte, 4096)

	for {

		if n, err := conn.Read(buffer); err != nil {
			checkError(err)
		} else {
			go handleProtoClient(buffer[0:n])
		}

	}

}

func handleProtoClient(data []byte) {

	m := new(Messages.NetworkMessage)

	err := proto.Unmarshal(data, m)
	checkError(err)

	if m.Request == true {
		if m.Type == 1 {
			log.Println("Unpacking message")

			m2 := new(Messages.Dabes)
			err := proto.Unmarshal(m.Payload, m2)
			checkError(err)

			log.Println("Messaged recieved:", m2.Message)
		}
	}

}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
}

func main() {
	Listen("0.0.0.0", 45678)

	fmt.Scanln()
}
