package kademlia

import (
	"log"

	"github.com/ArmedGuy/kadfs/message"
)

func (network *Network) registerMessageHandlers() {
	network.SetRequestHandler("FIND_NODE", func(sender *Contact, rpc *RPCMessage) {
		req := new(message.FindNodeRequest)
		rpc.GetMessageFromPayload(req)
		target := NewKademliaID(req.TargetID)

		contacts := network.kademlia.RoutingTable.FindClosestContacts(target, K)
		res := new(message.FindNodeResponse)
		for _, c := range contacts {
			res.Contacts = append(res.Contacts, &message.Contact{
				ID:      c.ID.String(),
				Address: c.Address,
			})
		}
		resRpc := rpc.GetResponse()
		resRpc.SetPayloadFromMessage(res)
		network.SendUDPPacket(sender, resRpc.GetBytes())
	})

	network.SetRequestHandler("FIND_VALUE", func(sender *Contact, rpc *RPCMessage) {
		log.Println("Find value")
	})

	network.SetRequestHandler("PING", func(sender *Contact, rpc *RPCMessage) {
		log.Println("Ping")
	})
}
