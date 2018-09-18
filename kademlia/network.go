package kademlia

import (
	"encoding/binary"
	"log"
	"net"

	"../message"
	"github.com/golang/protobuf/proto"
)

type Network struct {
	Me        *Contact
	Requests  map[int32]func(message.RPC, []byte)
	Responses map[int32]func(message.RPC, []byte)
}

func HandleError(err error) {
	log.Fatal(err)
}

func (network *Network) Listen(stateq chan StateTransition) {
	log.Println("[INFO] kademlia: Listening, accepting RPCs on", network.Me.Address)
	addr, _ := net.ResolveUDPAddr("udp", network.Me.Address)
	conn, _ := net.ListenUDP("udp", addr)
	header := make([]byte, 4)

	for {
		if read, err := conn.Read(header); err != nil {
			HandleError(err) // Handle error
		} else {
			if read != 4 {
				HandleError(err) // Handle error
			}

			//
			// Continue deserialization into a generic RPC message
			//
			messageLength := int32(binary.BigEndian.Uint32(header))
			rpcMessageBuf := make([]byte, messageLength)

			if _, err := conn.Read(rpcMessageBuf); err != nil {
				HandleError(err) // Handle error
			}

			rpcMessage := new(message.RPC)
			if err = proto.Unmarshal(rpcMessageBuf, rpcMessage); err != nil {
				HandleError(err) // Handle error
			}

			// We always read the payload as well
			payloadBuf := make([]byte, rpcMessage.Length)
			if _, err := conn.Read(payloadBuf); err != nil {
				HandleError(err) // Handle error
			}

			if rpcMessage.Request {
				if callback, ok := network.Requests[rpcMessage.MessageId]; ok {
					go callback(*rpcMessage, payloadBuf)
				}
			} else {
				if callback, ok := network.Responses[rpcMessage.MessageId]; ok {
					go callback(*rpcMessage, payloadBuf)
				}
			}
		}
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
