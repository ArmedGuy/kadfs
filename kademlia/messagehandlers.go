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
		resRPC := rpc.GetResponse()
		resRPC.SetPayloadFromMessage(res)
		network.Transport.SendRPCMessage(sender, resRPC)
	})

	//
	// Handle FIND_VALUE requests
	// If a node has the file, it should respond with the file.
	// If it doesn't have the file, it should respond with the K closest contacts to the file id
	//
	network.SetRequestHandler("FIND_VALUE", func(sender *Contact, rpc *RPCMessage) {
		// Get request message
		req := new(message.FindValueRequest)
		rpc.GetMessageFromPayload(req)

		// First, check if we have the file
		file, ok := network.kademlia.FileMemoryStore.Get(req.Hash)

		// Buuh, we do not have the file, respond wiht K closest nodes to file
		if !ok {

			key := NewKademliaID(req.Hash)
			contacts := network.kademlia.RoutingTable.FindClosestContacts(key, K)
			res := new(message.FindValueResponse)
			res.HasData = false
			for _, c := range contacts {
				res.Contacts = append(res.Contacts, &message.Contact{
					ID:      c.ID.String(),
					Address: c.Address,
				})
			}
			resRPC := rpc.GetResponse()
			resRPC.SetPayloadFromMessage(res)
			network.Transport.SendRPCMessage(sender, resRPC)

		} else {
			// Wohoo, we have the file!
			res := new(message.FindValueResponse)
			res.HasData = true
			res.Data = *file

			resRPC := rpc.GetResponse()
			resRPC.SetPayloadFromMessage(res)
			network.Transport.SendRPCMessage(sender, resRPC)
		}
	})

	network.SetRequestHandler("PING", func(sender *Contact, rpc *RPCMessage) {
		res := rpc.GetResponse()
		network.Transport.SendRPCMessage(sender, res)
	})

	network.SetRequestHandler("STORE", func(sender *Contact, rpc *RPCMessage) {
		req := new(message.SendDataMessage)
		rpc.GetMessageFromPayload(req)

		log.Printf("[INFO]: Stored file on node %v\n", network.kademlia.Network.GetLocalContact().ID)

		network.kademlia.FileMemoryStore.Put(req.Hash, req.Data, false)

		resRPC := rpc.GetResponse()
		network.Transport.SendRPCMessage(sender, resRPC)
	})
}
