package kademlia

import (
	"log"

	"github.com/ArmedGuy/kadfs/message"
	"github.com/golang/protobuf/proto"
)

func (network *Network) registerMessageHandlers() {
	network.SetRequestHandler("FIND_NODE", func(msg message.RPC, data []byte) {
		rpcPayload := new(message.FindNodeRequest)
		if err := proto.Unmarshal(data, rpcPayload); err != nil {
			log.Printf("[WARNING] messagehandler: Could not deserialize rpc payload, error: %v\n", err)
			return
		}
		target := NewKademliaID(rpcPayload.TargetID)
		contacts := network.kademlia.RoutingTable.FindClosestContacts(target, K)
		rpcResponse := new(message.FindNodeResponse)
		for _, c := range contacts {
			rpcResponse.Contacts = append(rpcResponse.Contacts, &message.Contact{
				ID:      c.ID.String(),
				Address: c.Address,
			})
		}
		m := network.CreateRPCMessage(contact, messageID, "FIND_NODE", 0, false)
	})

	network.SetRequestHandler("FIND_VALUE", func(msg message.RPC, data []byte) {
		log.Println("Find value")
	})

	network.SetRequestHandler("PING", func(msg message.RPC, data []byte) {
		log.Println("Ping")
	})
}
