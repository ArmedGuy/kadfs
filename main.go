package main

import (
	"github.com/ArmedGuy/kadfs/kademlia"
)

func main() {
	network := &kademlia.Network{}
	go network.Listen()
}
