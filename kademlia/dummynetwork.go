package kademlia

type DummyNetwork struct {
	Nodes []*Kademlia
}

func (network *DummyNetwork) Listen() {

}
func (network *DummyNetwork) SendPingMessage(contact *Contact) {

}
func (network *DummyNetwork) SendFindContactMessage(contact *Contact) {

}
func (network *DummyNetwork) SendFindDataMessage(hash string) {

}
func (network *DummyNetwork) SendStoreMessage(hash string, data []byte) {

}
