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

	fmt.Printf("Files: %v\n", datastore.files)

	time.Sleep(5 * time.Second)

	v := datastore.Get("/t/1")

	fmt.Printf("get data: %v\n", *v)

	time.Sleep(1 * time.Second)

	datastore.Delete("/t/1")
	fmt.Printf("Files: %v\n", datastore.files)

	x := datastore.GetKeysForReplicate()

	fmt.Printf("Keys that needs to be republished: %v\n", x)

	time.Sleep(5 * time.Second)
	datastore.DeleteExpiredData()
	fmt.Printf("Files left after expired deletion: %v\n", datastore.files)
}
