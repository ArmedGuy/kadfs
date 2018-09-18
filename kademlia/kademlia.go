package kademlia

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
		Me: &me,
	}
	return state
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// Look up the Alpha closest to target in our local routing table
	contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, Alpha)

	//
	// Do we want to create a channel to store all k closest nodes in here
	// and pass that channel to the SendFindContactMessage function.
	// All responses will then be added into this channel etc? hmm just thinking...
	//

	// Send a message to these Alpha nodes to learn about their k closest to target
	for _, element := range contacts {
		// Channel or just go routine?

		go kademlia.Network.SendFindContactMessage(&element)
	}

	// Do recursivly until there are no more contacts to query or no new contacts are discovered
	// This is done in a "reciever thread" where we update the temporary list
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
