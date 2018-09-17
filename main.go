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

	// start udp listener
	go state.Network.Listen()
	// Start webserver, blocking in new goroutine
	go s3.ConfigureAndListen(":8080")
}
