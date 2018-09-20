package kademlia

import (
	"crypto/sha1"
	"sync"
	"time"
)

const tReplicate = 3600
const tRepublish = 86400
const tExpire = 86400

type Datastore interface {
	Put(path string, data []byte)
	Get(path string) ([]byte, error)
	Delete(path string)
	GetKeysForRepublishing() map[*KademliaID][]byte
	DeleteExpiredData()
}

type InMemoryStore struct {
	files     map[*KademliaID][]byte
	replicate map[*KademliaID]time.Time
	expire    map[*KademliaID]time.Time
	mutex     *sync.Mutex
}

func (store *InMemoryStore) init() {
	store.files = make(map[*KademliaID][]byte)
	store.replicate = make(map[*KademliaID]time.Time)
	store.expire = make(map[*KademliaID]time.Time)
	store.mutex = &sync.Mutex{}
}

func PathHash(path string) *KademliaID {
	h := sha1.New()
	h.Write([]byte(path))
	return NewKademliaID(string(h.Sum(nil)))
}

func (store *InMemoryStore) Put(path string, data []byte) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	hash := PathHash(path)
	store.files[hash] = data
	store.expire[hash] = time.Now().Add(tExpire * time.Second)
	store.replicate[hash] = time.Now().Add(tExpire * time.Second)

}

func (store *InMemoryStore) Get(path string) ([]byte, bool) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	hash := PathHash(path)

	file, ok := store.files[hash]

	return file, ok
}

func (store *InMemoryStore) Delete(path string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	hash := PathHash(path)

	delete(store.files, hash)
	delete(store.expire, hash)
	delete(store.replicate, hash)
}

func (store *InMemoryStore) GetKeysForRepublishing() map[*KademliaID][]byte {

}

func (store *InMemoryStore) DeleteExpiredData() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for key, value := range store.expire {
		if time.Now().After(value) {
			delete(store.files, key)
			delete(store.expire, key)
			delete(store.replicate, key)
		}
	}
}
