package kademlia

import (
	"log"
	"time"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	Network      KademliaNetwork
}

func NewKademliaState(me Contact, network KademliaNetwork) *Kademlia {
	state := &Kademlia{}
	state.RoutingTable = NewRoutingTable(me, state)
	state.Network = network
	network.SetState(state)
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

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
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
