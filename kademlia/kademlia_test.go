package kademlia

import (
	"log"
	"testing"
	"time"
)

func examineRoutingTable(state *Kademlia) {
	local := state.Network.GetLocalContact()
	log.Println("----------------------------------------------------------------------------------")
	log.Printf("Viewing routing table for node %v\n", local)
	for i, c := range state.RoutingTable.FindClosestContacts(local.ID, 20) {
		log.Printf("%v: %v at distance %v\n", i, c, local.ID.CalcDistance(c.ID))
	}
	log.Println("----------------------------------------------------------------------------------")
}

func TestKademliaLookup(t *testing.T) {
	id := NewRandomKademliaID()
	me := NewContact(id, "localhost:8001") // TODO: change

	id2 := NewRandomKademliaID()
	me2 := NewContact(id2, "localhost:8002") // TODO: change

	id3 := NewRandomKademliaID()
	me3 := NewContact(id3, "localhost:8003") // TODO: change

	log.Printf("[INFO] kadfs: Creating new state for %v with ID %v\n", me.Address, me.ID)
	network := NewNetwork(&me)
	network2 := NewNetwork(&me2)
	network3 := NewNetwork(&me3)
	state := NewKademliaState(me, network)
	state2 := NewKademliaState(me2, network2)
	state3 := NewKademliaState(me3, network3)

	go state.Network.Listen()

	go state2.Network.Listen()
	go state3.Network.Listen()

	time.Sleep(4 * time.Second)
	state.Bootstrap(&me2)
	time.Sleep(1 * time.Second)
	examineRoutingTable(state)
	examineRoutingTable(state2)
	state3.Bootstrap(&me2)
	time.Sleep(5 * time.Second)
	examineRoutingTable(state)
	examineRoutingTable(state2)
	examineRoutingTable(state3)

}
