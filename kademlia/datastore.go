package kademlia

type Datastore interface {
	Put(path string, data []byte)
	Get(path string) ([]byte, error)
	Delete(path string)
}

type InMemoryStore struct {
	files map[*KademliaID][]byte
}

/*func (store *InMemoryStore) Put(path string, data []byte]) {

}*/

/*func (store *InMemoryStore) Get(path string) ([]byte, error) {

}*/

/*func (store *InMemoryStore) Delete(path string) {

}*/
