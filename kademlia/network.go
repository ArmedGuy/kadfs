package kademlia

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"

	"../message"
	"github.com/golang/protobuf/proto"
)

type Network struct {
	Me            *Contact
	NextMessageID int32
	LocalAddress  *net.UDPAddr
}

func (network *Network) Listen(stateq chan StateTransition) {
	log.Println("[INFO] kademlia: Listening, accepting RPCs on", network.Me.Address)
	addr, _ := net.ResolveUDPAddr("udp", network.Me.Address)
	conn, _ := net.ListenUDP("udp", addr)
	network.LocalAddress = addr
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

func (network *Network) NextID() int32 {
	x := network.NextMessageID
	network.NextMessageID++
	return x
}

func (network *Network) SendUDPPacket(contact *Contact, data []byte) {
	raddr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		log.Printf("[WARNING] network: %v\n", err)
		return
	}

	conn, err := net.DialUDP("udp", network.LocalAddress, raddr)
	if err != nil {
		log.Printf("[WARNING] network: %v\n", err)
		return
	}
	conn.WriteToUDP(data, raddr)

}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactRequest(contact *Contact) {
	// Build the message and send a request to the contact
	messageID := network.NextID()
	m := network.CreateRPCMessage(contact, messageID, "FindContact", 0, true)

	// build using bytes.buffer
	rpcData, err := proto.Marshal(m)
	if err != nil {
		log.Printf("[WARNING] network: %v\n", err)
		return
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(rpcData)))

	var b bytes.Buffer
	b.Write(header)
	b.Write(rpcData)

	network.SendUDPPacket(contact, b.Bytes())
}

/*
func (network *Network) SendFindContactResponse(contact *Contact, messageId int32) {
	// Build the message and send a request to the contact

	m := network.CreateRPCMessage(contact, messageId, "FindContact", x, false)

	data, err := proto.Marshal(m)

	if err == nil {
		//send the packet
	}
}
*/

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) CreateRPCMessage(contact *Contact, messageId int32, remoteProcedure string, length int32, request bool) *message.RPC {
	m := new(message.RPC)
	m.SenderId = network.Me.ID.String()
	m.ReceiverId = contact.ID.String()
	m.RemoteProcedure = remoteProcedure
	m.Length = length
	m.Request = request
	m.SenderAddress = network.Me.Address
	m.MessageId = messageId

	return m
}
