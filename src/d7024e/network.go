package d7024e

import (
	"log"
	"net"
	"os"
	"strconv"
)

type Network struct {
}

func Listen(ip string, port int) {
	s := ip + ":" + strconv.Itoa(port)

	addr, err := net.ResolveUDPAddr("udp", s)
	checkError(err)

	conn, err := net.ListenUDP("udp", addr)
	checkError(err)

	buffer := make([]byte, 1500)

	for {
		if _, err := conn.Read(buffer); err != nil {
			log.Fatal(err)
		} else {
			// spawn thread that unpacks the message

		}

	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
