package main

import (
	"log"

	"github.com/ArmedGuy/kadfs/kademlia"
	"github.com/ArmedGuy/kadfs/s3"
)

func main() {

	id := kademlia.NewRandomKademliaID()
	me := kademlia.NewContact(id, "localhost:8001") // TODO: change

	log.Printf("[INFO] kadfs: Creating new state for %v with ID %v\n", me.Address, me.ID)
	state := kademlia.NewKademliaState(me)

	go state.Network.Listen(state.Queue)
	go s3.ConfigureAndListen(":8080")

	go func() {
		state.Queue <- &kademlia.RoutingTableEntry{
			FoundID: kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"),
			Address: "localhost:8002",
		}
	}()

	for {
		msg := <-state.Queue
		log.Printf("[DEBUG] kadfs: got state transition event")
		msg.Transition(state)
	}
}
