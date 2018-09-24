package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/ArmedGuy/kadfs/kademlia"
	"github.com/ArmedGuy/kadfs/s3"
)

func examineRoutingTable(state *kademlia.Kademlia) {
	local := state.Network.GetLocalContact()
	log.Println("----------------------------------------------------------------------------------")
	log.Printf("Viewing routing table for node %v\n", local)
	for i, c := range state.RoutingTable.FindClosestContacts(local.ID, 20) {
		log.Printf("%v: %v at distance %v\n", i, c, local.ID.CalcDistance(c.ID))
	}
	log.Println("----------------------------------------------------------------------------------")
}

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
		examineRoutingTable(state)
		examineRoutingTable(state2)
		state3.Bootstrap(&me2)
		time.Sleep(5 * time.Second)
		examineRoutingTable(state)
		examineRoutingTable(state2)
		examineRoutingTable(state3)

		fileToStore := []byte{1, 2, 3, 4, 5, 1, 3, 3, 7}
		h1 := sha1.New()
		h1.Write([]byte("some/file/path.exe"))

		n := state.Store(hex.EncodeToString(h1.Sum(nil)), fileToStore)
		log.Printf("[LOG]: %v answered the store\n", n)

		time.Sleep(5 * time.Second)

		// Try to find some value
		file, ok := state.FindValue(hex.EncodeToString(h1.Sum(nil)))
		log.Printf("Found file returned %v. File content: %v\n", ok, file)
	}()

	fmt.Scanln()

}
