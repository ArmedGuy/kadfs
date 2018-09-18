package kademlia

import "github.com/ArmedGuy/kadfs/message"

const Alpha = 3
const k = 20

type Kademlia struct {
	Queue        chan StateTransition
	RoutingTable *RoutingTable
	Network      *Network
}

func NewKademliaState(me Contact) *Kademlia {
	state := &Kademlia{}
	state.Queue = make(chan StateTransition)
	state.RoutingTable = NewRoutingTable(me)
	state.Network = &Network{
		Me:            &me,
		NextMessageID: 0,
		Requests:      make(map[string]func(message.RPC, []byte)),
		Responses:     make(map[int32]func(message.RPC, []byte)),
	}
	return state
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// Look up the Alpha closest to target in our local routing table
	contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, Alpha)

	// Send a message to these Alpha nodes to learn about their k closest to target
	for _, element := range contacts {
		// Channel or just go routine?
		go kademlia.Network.SendFindContactRequest(&element)

	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
