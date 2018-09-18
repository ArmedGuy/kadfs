package kademlia

import (
	"log"
	"net"
)

type Network struct {
	Me *Contact
}

func (network *Network) Listen(stateq chan StateTransition) {
	log.Println("[INFO] kademlia: Listening, accepting RPCs on", network.Me.Address)
	addr, _ := net.ResolveUDPAddr("udp", network.Me.Address)
	conn, _ := net.ListenUDP("udp", addr)
	header := make([]byte, 4)

	for {
		if read, err := conn.Read(header); err != nil {
			// Error handle somehow
			continue
		} else {
			if read != 4 {
				continue
			}

		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// Build the message and send to the contact

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
