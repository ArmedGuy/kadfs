package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"sync"
	"time"
)

const tReplicate = 3600
const tRepublish = 86400
const tExpire = 86400

type Datastore interface {
	Put(hash string, data []byte, isOriginal bool, expire int32)
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
	expire    time.Time
	republish time.Time
	Data      *[]byte
	isOG      bool
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

func (store *InMemoryStore) Put(hash string, data []byte, isOriginal bool, expire int32) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	var republish time.Time

	if isOriginal {
		republish = time.Now().Add(tRepublish * time.Second)
	} else {
		republish = time.Now().Add(tReplicate * time.Second)
	}

	store.files[hash] = &File{
		Data:      &data,
		republish: republish,
		expire:    time.Now().Add(time.Duration(expire) * time.Second),
		isOG:      isOriginal,
	}
}

func (store *InMemoryStore) GetFileObject(hash string) (*File, bool) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	s := store.files
	file, ok := s[hash]

	if !ok {
		return nil, false
	}

	return file, ok
}

func (store *InMemoryStore) GetData(hash string) (*[]byte, bool) {
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

func (store *InMemoryStore) Delete(hash string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	delete(store.files, hash)
}

func (store *InMemoryStore) GetKeysAndValueForReplicate() map[string]*File {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	temp := make(map[string]*File)
	for key, value := range store.files {
		if time.Now().After(value.republish) && !value.isOG {
			temp[key] = store.files[key]
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
		if time.Now().After(value.expire) {
			delete(store.files, key)
		}
	}
}

func (store *InMemoryStore) Update(hash string, data []byte, isOG bool, expire, republish time.Time) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	s := store.files
	file, ok := s[hash]

	if ok {
		file.Data = &data
		file.isOG = isOG
		file.expire = expire
		file.republish = republish
	}

}
