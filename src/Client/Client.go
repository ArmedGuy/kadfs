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

func Send(ip string, port int, data []byte) {
	s := ip + ":" + strconv.Itoa(port)

	raddr, err := net.ResolveUDPAddr("udp", s)
	checkError(err)

	laddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	checkError(err)

	conn, err := net.ListenUDP("udp", laddr)
	checkError(err)

	n, err := conn.WriteToUDP(data, raddr)
	checkError(err)

	log.Println("Sent %s bytes", strconv.Itoa(n))

}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
}

func main() {

	m := new(Messages.Dabes)

	m.Message = "tell me your best joke? ya mum"

	temp, err := proto.Marshal(m)

	checkError(err)

	m2 := new(Messages.NetworkMessage)
	m2.Type = 1
	m2.Request = true
	m2.Context = 2
	m2.Payload = temp

	data, err := proto.Marshal(m2)
	checkError(err)

	m3 := new(Messages.Dabes)

	m3.Message = "asdsadsad"

	temp2, err := proto.Marshal(m3)

	checkError(err)

	m4 := new(Messages.NetworkMessage)
	m4.Type = 1
	m4.Request = true
	m4.Context = 2
	m4.Payload = temp2

	data1, err := proto.Marshal(m4)
	checkError(err)

	for {
		Send("127.0.0.1", 45678, data)
		fmt.Scanln()

		Send("127.0.0.1", 45678, data1)
		fmt.Scanln()

	}
}
