package kademlia

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"

	"github.com/ArmedGuy/kadfs/message"
	"github.com/golang/protobuf/proto"
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
	Me            *Contact
	NextMessageID int32
	Conn          *net.UDPConn
	Requests      map[string]func(message.RPC, []byte)
	Responses     map[int32]func(message.RPC, []byte)
}

func (network *Network) GetLocalContact() *Contact {
	return network.Me
}

func (network *Network) Listen() {
	log.Println("[INFO] kademlia: Listening, accepting RPCs on", network.Me.Address)
	addr, _ := net.ResolveUDPAddr("udp", network.Me.Address)
	conn, _ := net.ListenUDP("udp", addr)
	network.Conn = conn
	buf := make([]byte, 4096) // Come up with a reasonable size for this!
	header := make([]byte, 4)


	for {
		if read, err := conn.Read(buf); err != nil {
			log.Printf("[WARNING] network: Could not read header, error: %v\n", err)
			return
		} else {
			b := bytes.NewBuffer(buf)
			read, _ = b.Read(header)
			if read != 4 {
				log.Printf("[WARNING] network: Incorrect header size read, got: %v\n", read)
				return
			}

			//
			// Continue deserialization into a generic RPC message
			//
			messageLength := int32(binary.BigEndian.Uint32(header))
			rpcMessageBuf := make([]byte, messageLength)

			if _, err := b.Read(rpcMessageBuf); err != nil {
				log.Printf("[WARNING] network: Could not read into rpcbuf, error: %v\n", err)
				return
			}

			rpcMessage := new(message.RPC)
			if err = proto.Unmarshal(rpcMessageBuf, rpcMessage); err != nil {
				log.Printf("[WARNING] network: Could not deserialize rpc, error: %v\n", err)
				return
			}

			// We always read the payload as well
			payloadBuf := make([]byte, rpcMessage.Length)
			if _, err := b.Read(payloadBuf); err != nil {
				log.Printf("[WARNING] network: Could not read payload into buffer, error: %v\n", err)
				return
			}

			// Map request/responses to function based on message remoteProcedure/ID
			if rpcMessage.Request {
				if callback, ok := network.Requests[rpcMessage.RemoteProcedure]; ok {
					go callback(*rpcMessage, payloadBuf)
				} else {
					log.Printf("No request handler for %v\n", rpcMessage.RemoteProcedure)
				}
			} else {
				if callback, ok := network.Responses[rpcMessage.MessageId]; ok {
					go callback(*rpcMessage, payloadBuf)
				} else {
					log.Printf("No response handler for %v\n", rpcMessage.MessageId)
				}
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
		log.Printf("[WARNING] network1: %v\n", err)
		return
	}
	network.Conn.WriteToUDP(data, raddr)
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

	// Register message response mapping for this unique message ID
	network.Responses[messageID] = HandleFindContactResponse

	log.Printf("Sent a find contact packet")
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

func (network *Network) SendStoreMessage(hash string, data []byte) {
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
