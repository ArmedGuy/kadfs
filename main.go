package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ArmedGuy/kadfs/kademlia"
	"github.com/ArmedGuy/kadfs/s3"
)

func main() {

	id := kademlia.NewRandomKademliaID()
	me := kademlia.NewContact(id, "localhost:8001") // TODO: change

	id2 := kademlia.NewRandomKademliaID()
	me2 := kademlia.NewContact(id2, "localhost:8002") // TODO: change

	id3 := kademlia.NewRandomKademliaID()
	me3 := kademlia.NewContact(id3, "localhost:8003") // TODO: change

	log.Printf("[INFO] kadfs: Creating new state for %v with ID %v\n", me.Address, me.ID)
	network := kademlia.NewNetwork(&me)
	network2 := kademlia.NewNetwork(&me2)
	network3 := kademlia.NewNetwork(&me3)
	state := kademlia.NewKademliaState(me, network)
	state2 := kademlia.NewKademliaState(me2, network2)
	state3 := kademlia.NewKademliaState(me3, network3)

	go state.Network.Listen()
	go s3.ConfigureAndListen(":8080")

	go state2.Network.Listen()
	go state3.Network.Listen()

	go func() {
		time.Sleep(4 * time.Second)
		state.Bootstrap(&me2)
		time.Sleep(1 * time.Second)
		state2.Bootstrap(&me3)
		time.Sleep(5 * time.Second)
		state3.RoutingTable.FindClosestContacts(id3, 20)
	}()

	fmt.Scanln()

}
