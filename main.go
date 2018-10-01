package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/ArmedGuy/kadfs/kademlia"
)

// Massive workaround because docker does not like 127.0.0.1
func GetInternalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()

}

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
	var myID *kademlia.KademliaID

	var origin = flag.Bool("origin", false, "should node be a bootstrap node")
	var listen = flag.String("listen", "0.0.0.0:4000", "which ip:port to listen for kademlia on")

	var bootstrapID = flag.String("bootstrap-id", "", "ID of node to bootstrap towards")
	var bootstrapIP = flag.String("bootstrap-ip", "", "IP of node to bootstrap towards")

	flag.Parse()

	if *origin {
		myID = kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	} else {
		rand.Seed(time.Now().UnixNano())
		myID = kademlia.NewRandomKademliaID()
	}

	me := kademlia.NewContact(myID, *listen)
	myNetwork := kademlia.NewNetwork(&me)

	state := kademlia.NewKademliaState(me, myNetwork)

	// Starting all go routines
	go state.Network.Listen()
	go func() {
		for {
			// should probably be different go routines with different time for updates
			timer := time.NewTimer(150 * time.Second)
			<-timer.C
			log.Printf("[INFO] kadfs: Running republish, expire and replicate")
			go state.Replicate()
			go state.Republish()
			go state.Expire()
		}
	}()

	if *origin {
		log.Println("[INFO] kadfs: Running in origin mode, no bootstrap!")
	} else {
		log.Printf("[INFO] kadfs: Sleeping for 2 seconds to make sure the bootstrap node is up.")
		time.Sleep(2 * time.Second)

		id2 := kademlia.NewKademliaID(*bootstrapID)
		bootstrapNode := kademlia.NewContact(id2, *bootstrapIP) // TODO: change

		// Should probably retry the boostrap a few times if we fail
		state.Bootstrap(&bootstrapNode)
	}

	// Listen for user input here whenever that gets implemented
	// fmt.Scanln()

	for {
		time.Sleep(5 * time.Second)
		examineRoutingTable(state)
	}

}
