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
	nodes map[string]*Kademlia
}

type InternalRoutingTransport struct {
	global *kademliatestnetwork
	From   *Contact
}

func (trans *InternalRoutingTransport) SendRPCMessage(to *Contact, rpc *RPCMessage) {
	otherstate, ok := trans.global.nodes[rpc.Header.ReceiverId]
	if !ok {
		log.Printf("[ERROR] InternalRoutingTransport: No node in network matching ID")
		return
	}
	othernetwork, ok := otherstate.Network.(*Network)
	if !ok {
		log.Printf("[ERROR]: Other network broken, cannot do internal message routing")
	}
	otherstate.RoutingTable.AddContact(*trans.From)
	if rpc.Header.Request {
		callback, _ := othernetwork.Requests[rpc.Header.RemoteProcedure]
		callback(trans.From, rpc)
	} else {
		callback, _ := othernetwork.Responses[rpc.Header.MessageId]
		callback(trans.From, rpc)
		delete(othernetwork.Responses, rpc.Header.MessageId)
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

func createKademliaNode(offset int, global *kademliatestnetwork) *Kademlia {
	id := NewRandomKademliaID()
	me := NewContact(id, fmt.Sprintf("localhost:%v", 8000+offset))
	network := NewNetwork(&me)
	// Set internal transport solution so we dont need UDP ports
	network.Transport = &InternalRoutingTransport{global: global, From: &me}
	return NewKademliaState(me, network)
}

func createKademliaNetwork(count int) *kademliatestnetwork {
	testnet := new(kademliatestnetwork)
	testnet.nodes = make(map[string]*Kademlia)
	node1 := createKademliaNode(0, testnet)
	testnet.nodes[node1.Network.GetLocalContact().ID.String()] = node1
	for i := 1; i < count; i++ {
		node := createKademliaNode(i, testnet)
		testnet.nodes[node.Network.GetLocalContact().ID.String()] = node
		node.Bootstrap(node1.Network.GetLocalContact())
	}
	return testnet
}

func TestKademliaBootstrap(t *testing.T) {

	testnet := createKademliaNetwork(20)
	for k := range testnet.nodes {
		log.Printf("node %v", testnet.nodes[k].Network.GetLocalContact())
		examineRoutingTable(testnet.nodes[k])
	}
}

func TestKademliaEviction(t *testing.T) {
	// Create 20 nodes that all end up in same bucket
	// disable oldest node (easiest is to overwrite response handler for PING)
	// create 1 extra node and insert in routing table
	// the oldest node should be evicted.
}

func TestKademliaNoEviction(t *testing.T) {
	// create 20 nodes that all end up in the same bucket
	// create 1 extra node and insert into routing table
	// no change to bucket should be made
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
