package kademlia

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
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
