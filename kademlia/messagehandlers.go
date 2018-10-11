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
		file, ok := network.kademlia.FileMemoryStore.GetEntireFile(req.Hash)

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
			res.Data = *file.Data
			res.OriginalPublisher = &message.Contact{
				ID:      file.OriginalPublisher.ID.String(),
				Address: file.OriginalPublisher.Address,
			}

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

		// Create a original publisher contact from message info and store that in file struct (used for delete etc)
		originalPublisher := NewContact(NewKademliaID(req.OriginalPublisherID), req.OriginalPublisherAddr)
		network.kademlia.FileMemoryStore.Put(&originalPublisher, req.Hash, req.Data, false)

		resRPC := rpc.GetResponse()
		network.Transport.SendRPCMessage(sender, resRPC)
	})

	network.SetRequestHandler("DELETE", func(sender *Contact, rpc *RPCMessage) {
		req := new(message.DeleteValueRequest)
		rpc.GetMessageFromPayload(req)

		log.Printf("[INFO] MessageHandlers Delete: Got DELETE message for file %v on node %v\n", req.Hash, network.kademlia.Network.GetLocalContact().ID)

		fileWasDeleted := network.kademlia.FileMemoryStore.Delete(req.Hash)

		resRPC := rpc.GetResponse()

		// Here we need to create a DeletResponse and add as payload to resRPC
		res := new(message.DeleteValueResponse)
		if fileWasDeleted {
			// Respond with true
			res.Deleted = true
			log.Printf("[INFO] MessageHandlers Delete: File %v was successfully deleted from node %v\n", req.Hash, network.kademlia.Network.GetLocalContact().ID)
		} else {
			// Respond with false
			log.Printf("[INFO] MessageHandlers Delete: File %v could not be deleted from node %v\n", req.Hash, network.kademlia.Network.GetLocalContact().ID)
			res.Deleted = false
		}
		resRPC.SetPayloadFromMessage(res)
		network.Transport.SendRPCMessage(sender, resRPC)
	})
}
