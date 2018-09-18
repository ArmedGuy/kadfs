package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ArmedGuy/kadfs/kademlia"
	"github.com/ArmedGuy/kadfs/message"
	"github.com/ArmedGuy/kadfs/s3"
)

func main() {

	id := kademlia.NewRandomKademliaID()
	me := kademlia.NewContact(id, "localhost:8001") // TODO: change

	id2 := kademlia.NewRandomKademliaID()
	me2 := kademlia.NewContact(id2, "localhost:8002") // TODO: change

	log.Printf("[INFO] kadfs: Creating new state for %v with ID %v\n", me.Address, me.ID)
	network := kademlia.NewNetwork(&me)
	network2 := kademlia.NewNetwork(&me2)
	state := kademlia.NewKademliaState(me, network)
	state2 := kademlia.NewKademliaState(me2, network2)

	go state.Network.Listen()
	go s3.ConfigureAndListen(":8080")

	go state2.Network.Listen()

	state.Network.SetRequestHandler("FindContact", func(rpc message.RPC, data []byte) {
		log.Println("Finding all the contacts")
	})

	go func() {
		time.Sleep(4 * time.Second)
		state2.Network.SendFindContactMessage(&me, kademlia.NewRandomKademliaID(), make(chan *kademlia.LookupResponse))
	}()

	fmt.Scanln()

}
