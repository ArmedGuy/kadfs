package kademlia

import (
	"fmt"
	"sort"
	"sync"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
type Contact struct {
	ID          *KademliaID
	Address     string
	distance    *KademliaID
	isAvailable bool
	lock        *sync.RWMutex
}

// NewContact returns a new instance of a Contact
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil, true, &sync.RWMutex{}}
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

// String returns a simple string representation of a Contact
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

func (contact *Contact) IsAvailable() bool {
	contact.lock.RLock()
	res := contact.isAvailable
	contact.lock.RUnlock()
	return res
}

func (contact *Contact) SetAvailable(set bool) {
	contact.lock.Lock()
	contact.isAvailable = set
	contact.lock.Unlock()
}

// ContactCandidates definition
// stores an array of Contacts
type ContactCandidates struct {
	contacts []Contact
}

// Append an array of Contacts to the ContactCandidates
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

// GetContacts returns the first count number of Contacts
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.contacts[:count]
}

func (candidates *ContactCandidates) GetAvailableContacts(count int) []Contact {
	var availContacts []Contact
	for _, c := range candidates.contacts {
		if c.IsAvailable() {
			availContacts = append(availContacts, c)
		}
		if count--; count == 0 {
			break
		}
	}
	return availContacts
}

// Sort the Contacts in ContactCandidates
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the ContactCandidates
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}

/// Table for assisting with lookups. Similar to ContactCandidates but also holds queried state

type LookupCandidate struct {
	Contact *Contact
	Queried bool
}

type TemporaryLookupTable struct {
	candidates   []*LookupCandidate
	LookupTarget *KademliaID
}

// Append an array of Contacts to the TemporaryLookupTable, and check if it already exists
func (table *TemporaryLookupTable) Append(contacts []Contact) {
	var newCandidates []*LookupCandidate
	for _, c := range contacts {
		newCandidates = append(newCandidates, &LookupCandidate{Contact: &c, Queried: false})
	}
	table.candidates = append(table.candidates, newCandidates...)
}

// GetContacts returns the first count number of Contacts
func (table *TemporaryLookupTable) GetContacts(count int) []Contact {
	var contacts []Contact
	for _, c := range table.candidates {
		contacts = append(contacts, *c.Contact)
		if count--; count == 0 {
			break
		}
	}
	return contacts
}

func (table *TemporaryLookupTable) GetAvailableContacts(count int) []Contact {
	var availContacts []Contact
	for _, c := range table.candidates {
		if c.Contact.IsAvailable() {
			availContacts = append(availContacts, *c.Contact)
		}
		if count--; count == 0 {
			break
		}
	}
	return availContacts
}

func (table *TemporaryLookupTable) GetNewCandidates(count int) []*LookupCandidate {
	var availCandidates []*LookupCandidate
	for _, c := range table.candidates {
		if !c.Queried && c.Contact.IsAvailable() {
			availCandidates = append(availCandidates, c)
		}
		if count--; count == 0 {
			break
		}
	}
	return availCandidates
}

// Sort the Contacts in TemporaryLookupTable
func (table *TemporaryLookupTable) Sort() {
	sort.Sort(table)
}

// Len returns the length of the TemporaryLookupTable
func (table *TemporaryLookupTable) Len() int {
	return len(table.candidates)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (table *TemporaryLookupTable) Swap(i, j int) {
	table.candidates[i], table.candidates[j] = table.candidates[j], table.candidates[i]
}

// Less returns true if the the distance to target at index i is smaller than
// the distance to target at index j
func (table *TemporaryLookupTable) Less(i, j int) bool {
	return table.candidates[i].Contact.ID.CalcDistance(table.LookupTarget).Less(table.candidates[j].Contact.ID.CalcDistance(table.LookupTarget))
}
