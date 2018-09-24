package kademlia

import (
	"log"
	"time"
)

type Kademlia struct {
	RoutingTable    *RoutingTable
	Network         KademliaNetwork
	FileMemoryStore *InMemoryStore
}

func NewKademliaState(me Contact, network KademliaNetwork) *Kademlia {
	state := &Kademlia{}
	state.RoutingTable = NewRoutingTable(me, state)
	state.Network = network
	network.SetState(state)
	state.FileMemoryStore = NewInMemoryStore()
	return state
}

const K = 20
const alpha = 3

func (kademlia *Kademlia) Bootstrap(bootstrap *Contact) {
	log.Printf("[INFO] kademlia: Bootstrapping with contact %v\n", bootstrap)
	kademlia.RoutingTable.AddContact(*bootstrap)
	kademlia.FindNode(kademlia.Network.GetLocalContact().ID)
}

func (kademlia *Kademlia) FindNode(target *KademliaID) []Contact {
	contacts := kademlia.RoutingTable.FindClosestContacts(target, K)
	candidates := NewTemporaryLookupTable(kademlia.Network.GetLocalContact(), target)
	candidates.Append(contacts)
	candidates.Sort()
	// we call the last sendout a panic send
	panic := false
	changed := true
	// hold closest node seen, from our local routing table
	closest := candidates.GetAvailableContacts(1)[0].ID
	log.Printf("closest: %v\n", closest)
	for {
		// get alpha new candidates to send to
		sendto := candidates.GetNewCandidates(alpha)
		// special case, cant send to anybody, just return what I got now
		if len(sendto) == 0 {
			log.Println("out of sendtos")
			return candidates.GetAvailableContacts(K)
		}
		if !changed {
			log.Println("did not change")
			if panic {
				// panic already sent, return best list
				return candidates.GetAvailableContacts(K)
			}
			// "panic send"
			log.Println("sending panic")
			panic = true
			sendto = candidates.GetNewCandidates(K)
		} else {
			panic = false // reset panic if new closest found
		}
		closest = sendto[0].Contact.ID
		log.Printf("new closest: %v\n", closest)
		// create a shared channel for all our responses
		reschan := make(chan *LookupResponse)
		handled := len(sendto)
		for _, candidate := range sendto {
			// We update the lookup candidate to queried state.
			// this means that it wont end up in GetNewCandidates queries
			log.Printf("[DEBUG] kademlia: Sending message to %v\n", candidate.Contact)
			candidate.Queried = true
			//go kademlia.Network.SendFindNodeBlaBla(candidate.Contact, reschan)
			go kademlia.Network.SendFindNodeMessage(candidate.Contact, target, reschan)
		}
		for handled > 0 {
			// select response from channel or a timeout
			// because recv has timeouts, sends must also have a timeout (albeit higher)
			select {
			case response := <-reschan: // got a response, remove responder from sendto list
				tmp := sendto[:0]
				for _, c := range sendto {
					if !c.Contact.ID.Equals(response.From.ID) {
						tmp = append(tmp, c)
					}
				}
				log.Println("got response")
				for _, c := range response.Contacts {
					log.Printf("got contact %v\n", c)
				}
				sendto = tmp
				candidates.Append(response.Contacts)
				break
			case <-time.After(3 * time.Second):
				log.Println("timeout of response")
				break
			}
			handled--
		}
		// any node still in sendto list after all are handled are considered timed out
		for _, c := range sendto {
			c.Contact.SetAvailable(false)
		}
		candidates.Sort()
		// calculate if best candidates have changed or not
		newClosest := candidates.GetAvailableContacts(1)
		if len(newClosest) == 0 {
			log.Println("No available contacts!")
			return newClosest
		}
		newClosestID := newClosest[0].ID
		changed = false
		if newClosestID.CalcDistance(target).Less(closest.CalcDistance(target)) {
			changed = true
			closest = newClosestID
		}
	}
}

func (kademlia *Kademlia) FindValue(hash string) ([]byte, bool) { // Return File and a bool indicating error (file not found?)
	// Before we do anything, we check if we have the file
	log.Printf("[INFO] Kademlia FindValue: Searching for file in local storage\n")
	file, ok := kademlia.FileMemoryStore.Get(hash)
	log.Printf("[INFO] Kademlia FindValue: Local file status: %v, File: %v\n", ok, file)

	if ok {
		return *file, true
	}

	// Variables needed for the FindFile procedure
	panic := false
	changed := true

	// We start by converting the hashed file path to an kademlia ID
	key := HashToKademliaID(hash)

	// Get our K closest nodes to this key
	contacts := kademlia.RoutingTable.FindClosestContacts(key, K)
	candidates := NewTemporaryLookupTable(kademlia.Network.GetLocalContact(), key) // Append myself into this table
	candidates.Append(contacts)
	candidates.Sort()

	for {
		// Get alpha number of clients to send FIND_VALUE request to
		sendTo := candidates.GetNewCandidates(alpha)

		// No contacts available
		if len(sendTo) == 0 {
			log.Println("[INFO] kademlia FindValue: Found no contacts to send FindValue request to")

			// RETURN: What should be returned if no file is found?
			return nil, false
		}

		if !changed {
			log.Println("[INFO] kademlia FindValue: Could not find any closer nodes to the file")
			if panic {
				log.Println("[INFO] kademlia FindValue: Found no contacts to send FindValue request to")
				// Panic already sent

				// RETURN: What should be returned if no file is found?
				return nil, false
			}

			// Set panic
			log.Println("[INFO] kademlia FindValue: PANIC set")
			panic = true
			sendTo = candidates.GetNewCandidates(K)
		} else {
			panic = false // Closest contact changed, do not panic
		}

		closestNode := sendTo[0].Contact.ID
		log.Printf("[INFO] kademlia FindValue: Closest node is %v\n", closestNode)

		// Create a shared channel for responses
		responseChannel := make(chan *FindValueResponse)
		clientsToHandle := len(sendTo)

		// Query each candidate and update client query state
		for _, candidate := range sendTo {
			log.Printf("[INFO] kademlia FindValue: Sending message to %v\n", candidate.Contact)
			candidate.Queried = true
			go kademlia.Network.SendFindValueMessage(candidate.Contact, hash, responseChannel)
		}

		// Handle every client
		for clientsToHandle > 0 {
			select {
			case response := <-responseChannel: // got a response, remove responder from sendTo list

				// Did we get the file back?
				if response.HasFile {
					return *response.File.Data, true // WIN WIN, WE FOUND THE FILE!!!!!!!!!!!!
				} else {
					// We did not get any file back...
					// Append all contacts (exactly the same as in FindNode)
					tmp := sendTo[:0]
					for _, c := range sendTo {
						if !c.Contact.ID.Equals(response.From.ID) {
							tmp = append(tmp, c)
						}
					}
					log.Println("[INFO] Kademlia FindValue: got response")
					for _, c := range response.Contacts {
						log.Printf("[INFO] Kademlia FindValue: FindValuegot contact %v\n", c)
					}
					sendTo = tmp
					candidates.Append(response.Contacts)
					break
				}
			case <-time.After(3 * time.Second):
				log.Println("[INFO] Kademlia FindValue: timeout of response")
				break
			}
			clientsToHandle--
		}

		// All nodes still in the sendTo list have timed out
		for _, c := range sendTo {
			c.Contact.SetAvailable(false)
		}

		candidates.Sort()

		// calculate if best candidates have changed or not
		newClosest := candidates.GetAvailableContacts(1)
		if len(newClosest) == 0 {
			return nil, false // No available contacts left...
		}
		newClosestID := newClosest[0].ID
		changed = false
		if newClosestID.CalcDistance(key).Less(closestNode.CalcDistance(key)) {
			changed = true
			closestNode = newClosestID
		}
	}
}

func (kademlia *Kademlia) Store(hash string, data []byte) int {
	reschan := make(chan bool)
	closest := kademlia.FindNode(NewKademliaID(hash))

	kademlia.FileMemoryStore.Put(hash, data, true)

	for _, node := range closest {
		if node.ID != kademlia.Network.GetLocalContact().ID {
			n := node
			go kademlia.Network.SendStoreMessage(&n, hash, data, reschan)
		}
	}

	clientsToHandle := len(closest)

	responseAmount := 0

	for clientsToHandle > 0 {
		select {
		case <-reschan:
			responseAmount++
		case <-time.After(5 * time.Second):
			break
		}
		clientsToHandle--
	}
	return responseAmount

}

func (kademlia *Kademlia) Ping(contact *Contact) bool {
	reschan := make(chan bool)
	go kademlia.Network.SendPingMessage(contact, reschan)
	select {
	case <-reschan:
		return true
	case <-time.After(1 * time.Second):
		return false
	}
}

func (kademlia *Kademlia) Republish() {
	m := kademlia.FileMemoryStore.GetKeysAndValueForRepublish()

	for key, value := range m {
		_ = kademlia.Store(key, *value.Data)
	}
}

func (kademlia *Kademlia) Replicate() {
	keys := kademlia.FileMemoryStore.GetKeysForReplicate()

	for _, key := range keys {
		closest := kademlia.FindNode(NewKademliaID(key))

		iAmInClosest := false
		for _, contact := range closest {
			if kademlia.Network.GetLocalContact().ID.Equals(contact.ID) {
				iAmInClosest = true
			}
		}

		if iAmInClosest {
			// update time
			kademlia.FileMemoryStore.UpdateReplicateTime(key)
		} else {
			kademlia.FileMemoryStore.Delete(key)
		}

	}
}

func (kademlia *Kademlia) Expire() {
	kademlia.FileMemoryStore.DeleteExpiredData()
}
