package kademlia

import (
	"time"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	Network      KademliaNetwork
}

func NewKademliaState(me Contact, network KademliaNetwork) *Kademlia {
	state := &Kademlia{}
	state.RoutingTable = NewRoutingTable(me)
	state.Network = network
	return state
}

type LookupResponse struct {
	From     *Contact
	Contacts []Contact
}

const alpha = 3

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, 20)
	var candidates TemporaryLookupTable
	// load the lookup table with target ID. Is used to sort table with closest first
	candidates.LookupTarget = target.ID
	candidates.Append(contacts)
	candidates.Sort()
	// we call the last sendout a panic send
	panic := false
	changed := true
	// hold closest node seen, from our local routing table
	closest := candidates.GetAvailableContacts(1)[0].ID
	for {
		// get alpha new candidates to send to
		sendto := candidates.GetNewCandidates(alpha)
		if !changed {
			if panic {
				// panic already sent, return best list
				return candidates.GetAvailableContacts(20)
			}
			// "panic send"
			panic = true
			sendto = candidates.GetNewCandidates(20)
		} else {
			panic = false // reset panic if new closest found
		}
		closest = sendto[0].Contact.ID
		// create a shared channel for all our responses
		reschan := make(chan *LookupResponse)
		handled := len(sendto)
		for _, candidate := range sendto {
			// We update the lookup candidate to queried state.
			// this means that it wont end up in GetNewCandidates queries
			candidate.Queried = true
			//go kademlia.Network.SendFindNodeBlaBla(candidate.Contact, reschan)
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
				sendto = tmp
				candidates.Append(response.Contacts)
				break
			case <-time.After(3 * time.Second):
				break
			}
			handled--
		}
		// any node still in sendto list after all are handled are considered timed out
		// this updates the node-global state which means that we might evict this contact from a bucket
		for _, c := range sendto {
			c.Contact.SetAvailable(false)
		}
		candidates.Sort()
		// calculate if best candidates have changed or not
		newClosest := candidates.GetAvailableContacts(1)[0].ID
		changed = false
		if newClosest.CalcDistance(target.ID).Less(closest.CalcDistance(target.ID)) {
			changed = true
			closest = newClosest
		}
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
