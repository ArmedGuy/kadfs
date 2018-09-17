package kademlia

type StateTransition interface {
	Transition(*Kademlia)
}

type RoutingTableEntry struct {
	FoundID *KademliaID
	Address string
}

func (entry *RoutingTableEntry) Transition(state *Kademlia) {
	state.RoutingTable.AddContact(NewContact(entry.FoundID, entry.Address))

}
