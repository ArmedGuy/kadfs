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
	data      *[]byte
}

// Dno if this is needed or if we create all the stuff somewhere else
func (store *InMemoryStore) Init() {
	store.files = make(map[string]*File)
	store.mutex = &sync.Mutex{}
}

func PathHash(path string) string {
	h := sha1.New()
	h.Write([]byte(path))
	return NewKademliaID(hex.EncodeToString(h.Sum(nil))).String()
}

func (store *InMemoryStore) Put(path string, data []byte) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	hash := PathHash(path)

	store.files[hash] = &File{
		data:      &data,
		replicate: time.Now().Add(tReplicate * time.Second),
		expire:    time.Now().Add(tExpire * time.Second),
	}

}

func (store *InMemoryStore) Get(path string) *[]byte {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	hash := PathHash(path)

	s := store.files
	s1 := s[hash]
	file := s1.data

	return file
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
