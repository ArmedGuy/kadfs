package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"sync"
	"time"
)

// Real constants. Where do i use tRepublish and where do i use tReplicate?
// What is the difference?
const tReplicate = 3600
const tRepublish = 86400
const tExpire = 86400

type Datastore interface {
	Put(path string, data []byte)
	Get(path string) []byte
	Delete(path string)
	GetKeysForRepublishing() map[string]*File
	DeleteExpiredData()
}

type InMemoryStore struct {
	files map[string]*File
	mutex *sync.Mutex
}

type File struct {
	replicate time.Time
	expire    time.Time
	Data      *[]byte
}

func NewInMemoryStore() *InMemoryStore {
	inMemoryStore := &InMemoryStore{}
	inMemoryStore.files = make(map[string]*File)
	inMemoryStore.mutex = &sync.Mutex{}
	return inMemoryStore
}

func PathHash(path string) string {
	h := sha1.New()
	h.Write([]byte(path))
	return NewKademliaID(hex.EncodeToString(h.Sum(nil))).String()
}

func HashToKademliaID(hash string) *KademliaID {
	return NewKademliaID(hash)
}

func (store *InMemoryStore) Put(path string, data []byte) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	hash := PathHash(path)

	store.files[hash] = &File{
		Data:      &data,
		replicate: time.Now().Add(tReplicate * time.Second),
		expire:    time.Now().Add(tExpire * time.Second),
	}
}

func (store *InMemoryStore) Get(path string) (*[]byte, bool) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	hash := PathHash(path)

	s := store.files
	s1, ok := s[hash]

	// No file found, return errrrrr
	if !ok {
		return nil, false
	}

	file := s1.Data
	return file, true
}

func (store *InMemoryStore) Delete(path string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	hash := PathHash(path)
	delete(store.files, hash)
}

// Do we need the data too or only the keys?
func (store *InMemoryStore) GetKeysForRepublishing() map[string]*File {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	temp := make(map[string]*File)
	for key, value := range store.files {
		if time.Now().After(value.replicate) {
			temp[key] = store.files[key]
		}
	}

	return temp
}

func (store *InMemoryStore) DeleteExpiredData() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for key, value := range store.files {
		if time.Now().After(value.expire) {
			delete(store.files, key)
		}
	}
}
