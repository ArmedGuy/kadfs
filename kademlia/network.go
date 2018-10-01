package kademlia

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"

	"github.com/ArmedGuy/kadfs/message"
	"github.com/golang/protobuf/proto"
)

type KademliaNetwork interface {
	GetLocalContact() *Contact
	Listen()
	SendPingMessage(*Contact, chan bool)
	SendFindNodeMessage(*Contact, *KademliaID, chan *LookupResponse)
	SendFindValueMessage(*Contact, string, chan *FindValueResponse)
	SendStoreMessage(*Contact, string, []byte, chan bool)
	SetRequestHandler(string, func(*Contact, *RPCMessage))
	SetState(*Kademlia)
}

type KademliaNetworkTransport interface {
	SendRPCMessage(*Contact, *RPCMessage)
}

type UDPNetworkTransport struct {
	network *Network
}

func (trans *UDPNetworkTransport) SendRPCMessage(to *Contact, rpc *RPCMessage) {
	trans.network.SendUDPPacket(to, rpc.GetBytes())
}

type Network struct {
	Me            *Contact
	NextMessageID int32
	Conn          *net.UDPConn
	Requests      map[string]func(*Contact, *RPCMessage)
	Responses     map[int32]func(*Contact, *RPCMessage)
	kademlia      *Kademlia
	Transport     KademliaNetworkTransport
	lock          *sync.RWMutex
}

func NewNetwork(me *Contact) *Network {
	network := &Network{
		Me:            me,
		NextMessageID: 0,
		Requests:      make(map[string]func(*Contact, *RPCMessage)),
		Responses:     make(map[int32]func(*Contact, *RPCMessage)),
		lock:          new(sync.RWMutex),
	}
	network.Transport = &UDPNetworkTransport{network: network}
	network.registerMessageHandlers()
	return network

}

func (network *Network) GetLocalContact() *Contact {
	return network.Me
}

func (network *Network) SetState(state *Kademlia) {
	network.kademlia = state
}

func (network *Network) SetRequestHandler(rpc string, fn func(*Contact, *RPCMessage)) {
	network.Requests[rpc] = fn
}

func (network *Network) GetResponseHandler(messageID int32) (func(*Contact, *RPCMessage), bool) {

	network.lock.Lock()
	defer network.lock.Unlock()
	res, ok := network.Responses[messageID]
	delete(network.Responses, messageID)
	return res, ok
}

func (network *Network) Listen() {
	log.Println("[INFO] kademlia: Listening, accepting RPCs on", network.Me.Address)
	addr, _ := net.ResolveUDPAddr("udp", network.Me.Address)
	conn, _ := net.ListenUDP("udp", addr)
	network.Conn = conn
	buf := make([]byte, 4096) // Come up with a reasonable size for this!

	for {
		if _, caddr, err := conn.ReadFromUDP(buf); err != nil { // TODO: extract and use caddr
			log.Printf("[WARNING] network: Could not read header, error: %v\n", err)
			return
		} else {
			log.Printf("[INFO] network: got a packet PEPE\n")
			sender := caddr.String()
			rpc := network.NewRPCFromDatagram(buf)
			contact := NewContact(NewKademliaID(rpc.Header.SenderId), sender)
			go network.kademlia.RoutingTable.AddContact(contact)
			// Map request/responses to function based on message remoteProcedure/ID
			if rpc.Header.Request {
				if callback, ok := network.Requests[rpc.Header.RemoteProcedure]; ok {
					go callback(&contact, rpc)
				} else {
					log.Printf("[WARNING] network: No request handler for %v\n", rpc.Header.RemoteProcedure)
				}
			} else {
				if callback, ok := network.GetResponseHandler(rpc.Header.MessageId); ok {
					go callback(&contact, rpc)
				} else {
					log.Printf("[WARNING] network: No response handler for %v\n", rpc.Header.MessageId)
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
		log.Printf("[WARNING] network: Could not resolve addr, error: %v\n", err)
		return
	}
	network.Conn.WriteToUDP(data, raddr)
}

func (network *Network) SendPingMessage(contact *Contact, reschan chan bool) {
	rpc := network.NewRPC(contact, "PING")
	messageID := rpc.GetMessageId()
	network.lock.Lock()
	network.Responses[messageID] = func(sender *Contact, rpc *RPCMessage) {
		select {
		case reschan <- true:
			break
		case <-time.After(5 * time.Second):
			break
		}
	}
	network.lock.Unlock()
	network.Transport.SendRPCMessage(contact, rpc)
}

type LookupResponse struct {
	From     *Contact
	Contacts []Contact
}

type FindValueResponse struct {
	From     *Contact
	HasFile  bool
	File     File
	Contacts []Contact
}

func (network *Network) SendFindNodeMessage(contact *Contact, target *KademliaID, reschan chan *LookupResponse) {
	// Build the message and send a request to the contact
	rpc := network.NewRPC(contact, "FIND_NODE")
	messageID := rpc.GetMessageId()

	payload := new(message.FindNodeRequest)
	payload.TargetID = target.String()

	rpc.SetPayloadFromMessage(payload)

	// Register message response mapping for this unique message ID
	network.lock.Lock()
	network.Responses[messageID] = func(sender *Contact, rpc *RPCMessage) {
		// deserialize data, turn into list of contacts, and add to result struct
		// we then pass on result to result channel
		resData := new(message.FindNodeResponse)
		rpc.GetMessageFromPayload(resData)
		res := &LookupResponse{From: contact}
		for _, c := range resData.Contacts {
			res.Contacts = append(res.Contacts, NewContact(NewKademliaID(c.ID), c.Address))
		}
		select {
		case reschan <- res:
			break
		case <-time.After(5 * time.Second):
			break // nobody read our channel after 5 seconds, they must assumed we timed out
		}
	}
	network.lock.Unlock()
	network.Transport.SendRPCMessage(contact, rpc)
}

//
// Send FIND_VALUE message to a contact
//
func (network *Network) SendFindValueMessage(contact *Contact, hash string, reschan chan *FindValueResponse) {
	rpc := network.NewRPC(contact, "FIND_VALUE")
	messageID := rpc.GetMessageId()

	payload := new(message.FindValueRequest)
	payload.Hash = hash

	rpc.SetPayloadFromMessage(payload)

	// Register response mapping for this message id
	network.lock.Lock()
	network.Responses[messageID] = func(sender *Contact, rpc *RPCMessage) {

		// Get response message
		responseMessage := new(message.FindValueResponse)
		rpc.GetMessageFromPayload(responseMessage)

		// Build a FindValueResponse from this message and add to channel
		findValueResponse := &FindValueResponse{From: contact}

		// Wohoo, we got the file!
		if responseMessage.HasData {
			findValueResponse.HasFile = true
			receivedFile := new(File)
			receivedFile.Data = &responseMessage.Data
			findValueResponse.File = *receivedFile
		} else {
			// We did not receive any file... Only got a bunch of contacts...
			for _, c := range responseMessage.Contacts {
				findValueResponse.Contacts = append(findValueResponse.Contacts, NewContact(NewKademliaID(c.ID), c.Address))
			}
		}

		// Add findValueResponse to channel
		select {
		case reschan <- findValueResponse:
			break
		case <-time.After(5 * time.Second):
			break // nobody read our channel after 5 seconds, they must assumed we timed out
		}
	}

	// Unlock and send message
	network.lock.Unlock()
	network.Transport.SendRPCMessage(contact, rpc)
}

func (network *Network) SendStoreMessage(contact *Contact, hash string, data []byte, reschan chan bool) {
	rpc := network.NewRPC(contact, "STORE")
	messageID := rpc.GetMessageId()

	payload := new(message.SendDataMessage)
	payload.Data = data
	payload.Hash = hash

	rpc.SetPayloadFromMessage(payload)

	network.lock.Lock()
	network.Responses[messageID] = func(sender *Contact, rpc *RPCMessage) {
		select {
		case reschan <- true:
			break
		case <-time.After(5 * time.Second):
			break
		}
	}

	// Unlock and send message
	network.lock.Unlock()
	network.Transport.SendRPCMessage(contact, rpc)

}

type RPCMessage struct {
	Header  *message.RPC
	payload []byte
}

func (network *Network) NewRPC(contact *Contact, remoteProcedure string) *RPCMessage {
	header := new(message.RPC)
	header.SenderId = network.Me.ID.String()
	header.ReceiverId = contact.ID.String()
	header.RemoteProcedure = remoteProcedure
	header.MessageId = network.NextID()
	header.Request = true

	rpc := &RPCMessage{
		Header:  header,
		payload: make([]byte, 0),
	}
	return rpc
}

func (network *Network) NewRPCFromDatagram(buf []byte) *RPCMessage {
	b := bytes.NewBuffer(buf)
	header := make([]byte, 4)
	read, _ := b.Read(header)
	if read != 4 {
		log.Printf("[WARNING] network: Incorrect header size read, got: %v\n", read)
		return nil
	}

	//
	// Continue deserialization into a generic RPC message
	//
	messageLength := int32(binary.BigEndian.Uint32(header))
	rpcMessageBuf := make([]byte, messageLength)

	if _, err := b.Read(rpcMessageBuf); err != nil {
		log.Printf("[WARNING] network: Could not read into rpcbuf, error: %v\n", err)
		return nil
	}

	rpcMessage := new(message.RPC)
	if err := proto.Unmarshal(rpcMessageBuf, rpcMessage); err != nil {
		log.Printf("[WARNING] network: Could not deserialize rpc, error: %v\n", err)
		return nil
	}

	// We always read the payload as well
	payloadBuf := make([]byte, rpcMessage.Length)
	if _, err := b.Read(payloadBuf); err != nil {
		log.Printf("[WARNING] network: Could not read payload into buffer, error: %v\n", err)
		return nil
	}

	return &RPCMessage{
		Header:  rpcMessage,
		payload: payloadBuf,
	}
}

func (builder *RPCMessage) GetMessageId() int32 {
	return builder.Header.MessageId
}

func (builder *RPCMessage) HasPayload() bool {
	return len(builder.payload) != 0
}

func (builder *RPCMessage) SetPayload(data []byte) {
	builder.payload = data
}

func (builder *RPCMessage) SetPayloadFromMessage(pb proto.Message) {
	var err error
	builder.payload, err = proto.Marshal(pb)
	if err != nil {
		log.Printf("[WARNING] network: Could not serialize rpc payload, error: %v\n", err)
		return
	}
}

func (builder *RPCMessage) GetMessageFromPayload(pb proto.Message) {
	var err error
	err = proto.Unmarshal(builder.payload, pb)
	if err != nil {
		log.Printf("[WARNING] network: Could not deserialize rpc payload, error: %v\n", err)
		return
	}
}

func (builder *RPCMessage) GetResponse() *RPCMessage {
	header := new(message.RPC)
	header.MessageId = builder.Header.MessageId
	header.RemoteProcedure = builder.Header.RemoteProcedure
	header.SenderId, header.ReceiverId = builder.Header.ReceiverId, builder.Header.SenderId
	header.Request = false
	header.Length = 0
	return &RPCMessage{
		Header:  header,
		payload: make([]byte, 0),
	}
}

func (builder *RPCMessage) GetBytes() []byte {
	var b bytes.Buffer
	header := make([]byte, 4)

	builder.Header.Length = int32(len(builder.payload))
	data, err := proto.Marshal(builder.Header)
	if err != nil {
		log.Printf("[WARNING] network: Could not serialize RPC header, error: %v\n", err)
	}

	binary.BigEndian.PutUint32(header, uint32(len(data)))

	b.Write(header)
	b.Write(data)

	if builder.Header.Length > 0 {
		b.Write(builder.payload)
	}
	return b.Bytes()

}
