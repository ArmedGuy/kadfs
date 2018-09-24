package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"testing"
	"time"
)

type kademliatestnetwork struct {
	origin   *Kademlia
	nodelist []*Kademlia
	nodes    map[string]*Kademlia
}

func (global *kademliatestnetwork) addToNetwork(node *Kademlia) {
	global.nodelist = append(global.nodelist, node)
	global.nodes[node.Network.GetLocalContact().ID.String()] = node

}

type InternalRoutingTransport struct {
	global *kademliatestnetwork
	From   *Contact
}

func (trans *InternalRoutingTransport) SendRPCMessage(to *Contact, rpc *RPCMessage) {
	otherstate, ok := trans.global.nodes[rpc.Header.ReceiverId]
	if !ok {
		log.Printf("[ERROR] InternalRoutingTransport: No node in network matching ID %v", rpc.Header.ReceiverId)
		return
	}
	othernetwork, ok := otherstate.Network.(*Network)
	if !ok {
		log.Printf("[ERROR]: Other network broken, cannot do internal message routing")
	}
	go otherstate.RoutingTable.AddContact(*trans.From)
	if rpc.Header.Request {
		callback, _ := othernetwork.Requests[rpc.Header.RemoteProcedure]
		go callback(trans.From, rpc)
	} else {
		log.Println("Sending response")

		if callback, ok := othernetwork.GetResponseHandler(rpc.Header.MessageId); ok {
			go callback(trans.From, rpc)
		} else {
			log.Println("No response found")
		}
		log.Println("Sent response")
	}
}

func examineRoutingTable(state *Kademlia) {
	local := state.Network.GetLocalContact()
	log.Println("----------------------------------------------------------------------------------")
	log.Printf("Viewing routing table for node %v\n", local)
	for i, c := range state.RoutingTable.FindClosestContacts(local.ID, 20) {
		log.Printf("%v: %v at distance %v\n", i, c, local.ID.CalcDistance(c.ID))
	}
	log.Println("----------------------------------------------------------------------------------")
}

var nextID byte = 0

func contactInList(contact *Contact, contacts []Contact) bool {
	for _, c := range contacts {
		if c.ID.Equals(contact.ID) {
			return true
		}
	}
	return false
}

func nextKademliaID() *KademliaID {
	defer func() { nextID++ }()
	str := fmt.Sprintf("000000000000000000000000000F0000000000%02X", nextID)

	return NewKademliaID(str)
}

func createKademliaNode(id *KademliaID, offset int, global *kademliatestnetwork) *Kademlia {
	me := NewContact(id, fmt.Sprintf("localhost:%v", 8000+offset))
	network := NewNetwork(&me)
	// Set internal transport solution so we dont need UDP ports
	network.Transport = &InternalRoutingTransport{global: global, From: &me}
	return NewKademliaState(me, network)
}

func createKademliaNetwork(count int, spread bool) *kademliatestnetwork {
	testnet := new(kademliatestnetwork)
	testnet.nodes = make(map[string]*Kademlia)
	var id1 *KademliaID
	if spread {
		id1 = NewRandomKademliaID()
	} else {
		id1 = NewKademliaID("0000000000000000000000000000000000000000")
	}
	node1 := createKademliaNode(id1, 0, testnet)
	testnet.origin = node1
	testnet.addToNetwork(node1)
	for i := 1; i < count; i++ {
		var id *KademliaID
		if spread {
			id = NewRandomKademliaID()
		} else {
			id = nextKademliaID()
		}
		node := createKademliaNode(id, i, testnet)
		testnet.addToNetwork(node)
		log.Printf("[DEBUG] kademlia_test: Bootstrapping node %v", node.Network.GetLocalContact().ID.String())
		node.Bootstrap(node1.Network.GetLocalContact())
	}
	return testnet
}

func TestKademliaBootstrap(t *testing.T) {

	// Create network with 21 nodes (1 node has 20 contacts)
	testnet := createKademliaNetwork(21, false)
	for _, n := range testnet.nodelist {
		log.Printf("node %v", n.Network.GetLocalContact())
		examineRoutingTable(n)
	}
}

func TestKademliaEviction(t *testing.T) {
	// Create 21 nodes that all end up in same bucket (1 node knows 20 other nodes, which is max bucket size)
	// disable all nodes (easiest is to overwrite response handler for PING)
	// create 1 extra node and insert in routing table
	// the oldest node should be evicted. (not checked yet)

	nodes := 21

	testnet := createKademliaNetwork(nodes, false)
	pinged := false
	for i := 1; i < nodes; i++ {
		testnet.nodelist[i].Network.SetRequestHandler("PING", func(sender *Contact, rpc *RPCMessage) {
			// Attempting to ping me
			log.Print("Attempting to ping downed node")
			pinged = true
		})
	}

	time.Sleep(1 * time.Second)
	addmeid := nextKademliaID()
	addme := createKademliaNode(addmeid, 50, testnet)
	testnet.addToNetwork(addme)
	addme.Bootstrap(testnet.origin.Network.GetLocalContact())
	time.Sleep(10 * time.Second)

	closest := testnet.origin.RoutingTable.FindClosestContacts(testnet.origin.Network.GetLocalContact().ID, 30)
	if !pinged {
		log.Fatal("No pings were sent in eviction case!")
	}
	if !contactInList(addme.Network.GetLocalContact(), closest) {
		log.Fatal("Not in closest 30 contacts! Should be in a bucket")
	}

}

func TestKademliaNoEviction(t *testing.T) {
	// Create 21 nodes that all end up in same bucket (1 node knows 20 other nodes, which is max bucket size)
	// Dont disable any nodes
	// create 1 extra node and insert in routing table
	// New node should not exist in bucket

	nodes := 21

	testnet := createKademliaNetwork(nodes, false)

	time.Sleep(1 * time.Second)
	addmeid := nextKademliaID()
	addme := createKademliaNode(addmeid, 50, testnet)
	testnet.addToNetwork(addme)
	addme.Bootstrap(testnet.origin.Network.GetLocalContact())
	time.Sleep(10 * time.Second)

	closest := testnet.origin.RoutingTable.FindClosestContacts(testnet.origin.Network.GetLocalContact().ID, 30)
	if contactInList(addme.Network.GetLocalContact(), closest) {
		log.Fatal("In closest 30 contacts! Should not be in buckets")
	}

}

func TestKademliaFindNodePanic(t *testing.T) {
	// Create enough nodes to trigger a panic during lookup
	// Panic should find 1 extra node after panic is done
}

func TestKademliaFindNodeTimeouts(t *testing.T) {
	// Create 20 nodes, and disable a few of them
	// FindNode should only return (20 - disabled) nodes
}

func TestKademliaFindValue(t *testing.T) {
	testnet := createKademliaNetwork(30)

	// Create file to store
	hash1 := sha1.New()
	hash1.Write([]byte("some/file/path/file.ext"))
	fileHashString := hex.EncodeToString(hash1.Sum(nil))
	fileContent := []byte{1, 2, 3, 4, 5, 1, 3, 3, 7}

	// Send store request to some node
	//n := testnet[0]
	//log.Printf("[LOG]: %v answered the store\n", n)

	// Wait for propagation
	time.Sleep(5 * time.Second)

	// Try to find some value
	//file, ok := state2.FindValue(hex.EncodeToString(h1.Sum(nil)))
	//log.Printf("Found file returned %v. File content: %v\n", ok, file)
}
