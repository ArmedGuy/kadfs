package kademlia

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

type KademliaNetwork interface {
	GetLocalContact() *Contact
	Listen()
	SendPingMessage(*Contact)
	SendFindContactMessage(*Contact)
	SendFindDataMessage(string)
	SendStoreMessage(string, []byte)
}

type Network struct {
	Me *Contact
}

func (network *Network) GetLocalContact() *Contact {
	return network.Me
}

func (network *Network) Listen() {
	log.Println("[INFO] kademlia: Listening, accepting RPCs on", network.Me.Address)
	addr, _ := net.ResolveUDPAddr("udp", network.Me.Address)
	conn, _ := net.ListenUDP("udp", addr)
	sizeOf := make([]byte, 4)

	for {
		// read int32
		readSize, err := conn.Read(sizeOf)
		if err != nil {
			continue
		}
		if readSize != 4 {
			continue
		}
		var size int
		_ = binary.Read(bytes.NewReader(sizeOf), binary.BigEndian, &size)
		// read header
		header := make([]byte, size)
		readHeader, err := conn.Read(header)
		if readHeader != size {
			continue
		}
		// deserialize header, get length, and read length

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

func (network *Network) SendStoreMessage(hash string, data []byte) {
	// TODO
}
