package kademlia

import (
	"fmt"
	"testing"
	"time"
)

func TestDatastore(t *testing.T) {
	datastore := InMemoryStore{}

	datastore.Init()

	datastore.Put("/t/1", []byte{1, 1})
	datastore.Put("/t/2", []byte{2, 2})
	datastore.Put("/t/3", []byte{3, 3})

	fmt.Printf("Files: %v\nExpire: %v\nReplicate: %v\n", datastore.files, datastore.expire, datastore.replicate)
	time.Sleep(5 * time.Second)

	datastore.Delete("/t/1")
	fmt.Printf("Files: %v\nExpire: %v\nReplicate: %v\n", datastore.files, datastore.expire, datastore.replicate)

	x := datastore.GetKeysForRepublishing()

	fmt.Printf("keys: %v", x)
}
