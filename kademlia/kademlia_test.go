package kademlia

import (
	"bytes"
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

		if callback, ok := othernetwork.GetResponseHandler(rpc.Header.MessageId); ok {
			go callback(trans.From, rpc)
		} else {
			log.Println("No response found")
		}
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
	me := NewContact(id, fmt.Sprintf("10.0.0.%v:%v", offset, 8000+offset))
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

//
// Unit test for FIND_VALUE RPC
// Should always be successful!!!
//
func TestKademliaStoreAndFindValue(t *testing.T) {
	testnet := createKademliaNetwork(100, true)

	// Use this node to send STORE RPC to
	firstNode := testnet.nodelist[0]

	// Use this node to send FIND_VALUE RPC to
	otherNode := testnet.nodelist[75]

	// Create file to store
	hash1 := sha1.New()
	hash1.Write([]byte("some/file/path/file.ext"))
	fileHashString := hex.EncodeToString(hash1.Sum(nil))
	fileContent := []byte{1, 2, 3, 4, 5, 1, 3, 3, 7}

	// Store file
	n := firstNode.Store(fileHashString, fileContent, true, tExpire)
	log.Printf("Stored file %v with content %v on %v number of responding nodes\n", fileHashString, fileContent, n)

	if n > K {
		log.Fatal("ERROR TestKademliaStoreAndFindValue: File stored on too many nodes")
	}

	// Sleep for propagation
	fmt.Println("Sleeping for 5 secs")
	time.Sleep(5 * time.Second)

	// Do a FIND_VALUE RPC
	file, ok := otherNode.FindValue(fileHashString)
	log.Printf("FIND_VALUE returned %v with file content: %v\n", ok, file)

	if !ok {
		log.Fatal("ERROR TestKademliaStoreAndFindValue: FIND_VALUE did not return ok for file with id " + fileHashString)
	}
}

//
// Unit test for FIND_VALUE RCP
// Should never find the file!
//
func TestKademliaNoStoreAndFindValue(t *testing.T) {
	testnet := createKademliaNetwork(30, true)

	// Use this node to send FIND_VALUE RPC to
	otherNode := testnet.nodelist[15]

	// Create file to store
	// We never store this file
	hash1 := sha1.New()
	hash1.Write([]byte("some/file/path/file.ext"))
	fileHashString := hex.EncodeToString(hash1.Sum(nil))

	// Do a FIND_VALUE RPC
	file, ok := otherNode.FindValue(fileHashString)
	log.Printf("FIND_VALUE returned %v with file content: %v\n", ok, file)

	if ok {
		log.Fatal("ERROR TestKademliaNoStoreAndFindValue: FIND_VALUE returned ok for a file with id " + fileHashString + " that has never been stored")
	}
}

//
// Unit test for storing one file
// and then try to find another one that isn't stored on the network
//
func TestKademliaStoreAndFindOtherValue(t *testing.T) {
	testnet := createKademliaNetwork(100, true)

	// Use this node to send STORE RPC to
	firstNode := testnet.nodelist[32]

	// Use this node to send FIND_VALUE RPC to
	otherNode := testnet.nodelist[75]

	// Create file to store
	hash1 := sha1.New()
	hash1.Write([]byte("some/file/path/file.ext"))
	fileHashString := hex.EncodeToString(hash1.Sum(nil))
	fileContent := []byte{1, 2, 3, 4, 5, 1, 3, 3, 7}

	// Store file
	n := firstNode.Store(fileHashString, fileContent, true, tExpire)
	log.Printf("Stored file %v with content %v on %v number of responding nodes\n", fileHashString, fileContent, n)

	if n > K {
		log.Fatal("ERROR TestKademliaStoreAndFindOtherValue: File stored on too many nodes")
	}

	// Sleep for propagation
	fmt.Println("Sleeping for 5 secs")
	time.Sleep(5 * time.Second)

	// Find this id
	hash2 := sha1.New()
	hash2.Write([]byte("some/file2/path/file.ext"))
	fileHashString2 := hex.EncodeToString(hash2.Sum(nil))

	// Do a FIND_VALUE RPC
	file, ok := otherNode.FindValue(fileHashString2)
	log.Printf("FIND_VALUE returned %v with file content: %v\n", ok, file)

	if ok {
		log.Fatal("ERROR TestKademliaStoreAndFindOtherValue: FIND_VALUE returned ok for file with id " + fileHashString2)
	}
}

//
// Unit test for testing overwriting data
//
func TestKademliaStoreOverwrite(t *testing.T) {
	testnet := createKademliaNetwork(100, true)

	// Use this node to send STORE RPC to
	firstNode := testnet.nodelist[32]

	// Use this node to send FIND_VALUE RPC to
	otherNode := testnet.nodelist[75]

	// Create file to store
	hash1 := sha1.New()
	hash1.Write([]byte("some/file/path/file.ext"))
	fileHashString := hex.EncodeToString(hash1.Sum(nil))
	fileContent1 := []byte{1, 2, 3, 4, 5, 1, 3, 3, 7}

	// Store file
	n := firstNode.Store(fileHashString, fileContent1, true, tExpire)
	log.Printf("Stored file %v with content %v on %v number of responding nodes\n", fileHashString, fileContent1, n)

	if n > K {
		log.Fatal("ERROR TestKademliaStoreOverwrite: File stored on too many nodes")
	}

	// Sleep for propagation
	fmt.Println("Sleeping for 5 secs")
	time.Sleep(5 * time.Second)

	// Create a new file with the same path but different content
	fileContent2 := []byte{7, 3, 2, 1, 0}

	// Store second file
	n = firstNode.Store(fileHashString, fileContent2, true, tExpire)
	log.Printf("Stored file %v with content %v on %v number of responding nodes\n", fileHashString, fileContent2, n)

	if n > K {
		log.Fatal("ERROR TestKademliaStoreOverwrite: File stored on too many nodes")
	}

	// Sleep for propagation
	fmt.Println("Sleeping for 5 secs")
	time.Sleep(5 * time.Second)

	// Do a FIND_VALUE RPC
	file, ok := otherNode.FindValue(fileHashString)
	log.Printf("FIND_VALUE returned %v with file content: %v\n", ok, file)

	if !ok {
		log.Fatal("ERROR TestKademliaStoreOverwrite: FIND_VALUE returned ok for file with id " + fileHashString)
	}

	if bytes.Compare(file, fileContent2) != 0 {
		log.Fatal("ERROR TestKademliaStoreOverwrite: Received wrong file in FIND_VALUE request")
	}
}

func TestKademliaFindNodePanic(t *testing.T) {
	// Create enough nodes to trigger a panic during lookup
	// Panic should find 1 extra node after panic is done
	testnet := createKademliaNetwork(1, false)
	nearId := nextKademliaID()
	for i := 0; i < 3; i++ {
		id := nextKademliaID()
		node := createKademliaNode(id, 1+i, testnet)
		testnet.addToNetwork(node)
		testnet.origin.RoutingTable.AddContact(*node.Network.GetLocalContact())

	}

	panicId := nextKademliaID()
	panicNode := createKademliaNode(panicId, 4, testnet)
	testnet.addToNetwork(panicNode)
	testnet.origin.RoutingTable.AddContact(*panicNode.Network.GetLocalContact())

	nearNode := createKademliaNode(nearId, 5, testnet)
	testnet.addToNetwork(nearNode)
	panicNode.RoutingTable.AddContact(*nearNode.Network.GetLocalContact())

	tryToFind := false
	nearNode.Network.SetRequestHandler("FIND_NODE", func(sender *Contact, rpc *RPCMessage) {
		tryToFind = true
		// this will timeout the response, but atleast we got the panic
	})

	testnet.origin.FindNode(nearId)
	time.Sleep(3)

	if !tryToFind {
		log.Fatal("Did not attempt to find node which was 4th in list. No panic sent")
	}

}

func TestKademliaFindNodeTimeouts(t *testing.T) {
	// Create 20 nodes, and disable a few of them
	// FindNode should only return (20 - disabled) nodes
}

func TestKademliaExpireTimer(t *testing.T) {
	testnet := createKademliaNetwork(100, true)
	firstNode := testnet.nodelist[0]

	// Create file to store
	hash1 := sha1.New()
	hash1.Write([]byte("some/file/path/file.ext"))
	fileHashString := hex.EncodeToString(hash1.Sum(nil))
	fileContent1 := []byte{1, 2, 3, 4, 5, 1, 3, 3, 7}

	// Store file with expire happening in 1 sec
	n := firstNode.Store(fileHashString, fileContent1, true, 0)

	log.Printf("Stored file %v with content %v on %v number of responding nodes\n", fileHashString, fileContent1, n)

	if n > K {
		log.Fatal("ERROR TestKademliaStoreOverwrite: File stored on too many nodes")
	}

	// Sleep for propagation
	fmt.Println("Sleeping for 1 secs")
	time.Sleep(1 * time.Second)

	// Should only expire the file on this node (thus if running findValue it should find a value somewhere in the network)
	firstNode.Expire()

	file, ok := firstNode.FileMemoryStore.GetData(fileHashString)
	log.Printf("GetData returned %v with file content: %v\n", ok, file)

	if ok || file != nil {
		log.Fatal("Found a expired piece of data on firstNode")
	}

}

func TestKademliaRepublishTimer(t *testing.T) {
	testnet := createKademliaNetwork(20, true)
	firstNode := testnet.nodelist[15]

	// Create file to store
	hash1 := sha1.New()
	hash1.Write([]byte("some/file/path/file.ext"))
	fileHashString := hex.EncodeToString(hash1.Sum(nil))
	fileContent1 := []byte{1, 2, 3, 4, 5, 1, 3, 3, 7}

	// Store file
	n := firstNode.Store(fileHashString, fileContent1, true, 0)

	log.Printf("Stored file %v with content %v on %v number of responding nodes\n", fileHashString, fileContent1, n)

	if n > K {
		log.Fatal("ERROR TestKademliaStoreOverwrite: File stored on too many nodes")
	}

	// Sleep for propagation
	fmt.Println("Sleeping for 1 secs")
	time.Sleep(1 * time.Second)

	secondNode := testnet.nodelist[1]

	time.Sleep(1 * time.Second)

	firstNode.FileMemoryStore.Update(fileHashString, fileContent1, true, time.Now(), time.Now())
	secondNode.FileMemoryStore.Update(fileHashString, fileContent1, false, time.Now(), time.Now())

	firstNode.Republish()

	time.Sleep(1 * time.Second)

	file1, _ := firstNode.FileMemoryStore.GetFileObject(fileHashString)
	file2, _ := secondNode.FileMemoryStore.GetFileObject(fileHashString)

	if file1.republish != file2.republish {
		log.Fatalf("Node 1 and Node 2 have different republish times. Node1: %v, Node2: %v", file1.republish, file2.republish)
	}
}

func TestKademliaReplicateTimer(t *testing.T) {
	testnet := createKademliaNetwork(20, true)
	firstNode := testnet.nodelist[15]

	// Create file to store
	hash1 := sha1.New()
	hash1.Write([]byte("some/file/path/file.ext"))
	fileHashString := hex.EncodeToString(hash1.Sum(nil))
	fileContent1 := []byte{1, 2, 3, 4, 5, 1, 3, 3, 7}

	// Store file
	n := firstNode.Store(fileHashString, fileContent1, true, 0)

	log.Printf("Stored file %v with content %v on %v number of responding nodes\n", fileHashString, fileContent1, n)

	if n > K {
		log.Fatal("ERROR TestKademliaStoreOverwrite: File stored on too many nodes")
	}

	// Sleep for propagation
	fmt.Println("Sleeping for 1 secs")
	time.Sleep(1 * time.Second)

	secondNode := testnet.nodelist[1]

	time.Sleep(1 * time.Second)

	firstNode.FileMemoryStore.Update(fileHashString, fileContent1, true, time.Now(), time.Now())
	secondNode.FileMemoryStore.Update(fileHashString, fileContent1, false, time.Now(), time.Now())

	secondNode.Replicate()

	time.Sleep(1 * time.Second)

	file1, _ := firstNode.FileMemoryStore.GetFileObject(fileHashString)
	file2, _ := secondNode.FileMemoryStore.GetFileObject(fileHashString)

	if file1.republish != file2.republish {
		log.Fatalf("Node 1 and Node 2 have different republish times. Node1: %v, Node2: %v", file1.republish, file2.republish)
	}

}
