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
	Put(hash string, data []byte)
	Get(hash string) []byte
	Delete(hash string)
	GetKeysForReplicate() []string
	GetKeysAndValueForRepublish() map[string]*File
	DeleteExpiredData()
}

type InMemoryStore struct {
	files map[string]*File
	mutex *sync.Mutex
}

type File struct {
	OriginalPublisher *Contact
	replicate         time.Time
	expire            time.Time
	republish         time.Time
	Data              *[]byte
	isOG              bool
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

func (store *InMemoryStore) Put(originalPublisher *Contact, hash string, data []byte, isOriginal bool) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.files[hash] = &File{
		OriginalPublisher: originalPublisher,
		Data:              &data,
		republish:         time.Now().Add(tRepublish * time.Second),
		replicate:         time.Now().Add(tReplicate * time.Second),
		expire:            time.Now().Add(tExpire * time.Second),
		isOG:              isOriginal,
	}
}

func (store *InMemoryStore) Get(hash string) (*[]byte, bool) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	s := store.files
	s1, ok := s[hash]

	// No file found, return errrrrr
	if !ok {
		return nil, false
	}

	file := s1.Data
	return file, true
}

func (store *InMemoryStore) GetEntireFile(hash string) (*File, bool) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	s := store.files
	s1, ok := s[hash]

	// No file found, return errrrrr
	if !ok {
		return nil, false
	}

	return s1, true
}

func (store *InMemoryStore) Delete(hash string) bool {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	// First, check if we even have the file
	s := store.files
	_, ok := s[hash]

	// Delete if we have the file, else return false
	if !ok {
		return false
	}

	delete(store.files, hash)
	return true
}

func (store *InMemoryStore) GetKeysForReplicate() []string {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	temp := make([]string, 0)
	for key, value := range store.files {
		if time.Now().After(value.replicate) && !value.isOG {
			_ = append(temp, key)
		}
	}

	return temp
}

func (store *InMemoryStore) GetKeysAndValueForRepublish() map[string]*File {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	temp := make(map[string]*File)
	for key, value := range store.files {
		if time.Now().After(value.republish) && value.isOG {
			temp[key] = store.files[key]
		}
	}

	return temp
}

func (store *InMemoryStore) DeleteExpiredData() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for key, value := range store.files {
		if time.Now().After(value.expire) && !value.isOG {
			delete(store.files, key)
		}
	}
}

func (store *InMemoryStore) UpdateReplicateTime(hash string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.files[hash].replicate = time.Now().Add(tReplicate * time.Second)
}
