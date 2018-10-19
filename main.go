package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ArmedGuy/kadfs/s3"

	"github.com/ArmedGuy/kadfs/kademlia"

	"github.com/hashicorp/consul/api"
)

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalUnicast() && !ipnet.IP.IsLinkLocalMulticast() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func examineRoutingTable(state *kademlia.Kademlia) {
	local := state.Network.GetLocalContact()
	log.Println("----------------------------------------------------------------------------------")
	log.Printf("Viewing routing table for node %v\n", local)
	for i, c := range state.RoutingTable.FindClosestContacts(local.ID, 20) {
		log.Printf("%v: %v at distance %v\n", i, c, local.ID.CalcDistance(c.ID))
	}
	log.Println("----------------------------------------------------------------------------------")
}

func getRoutingTable(state *kademlia.Kademlia) string {
	local := state.Network.GetLocalContact()
	var b strings.Builder

	b.WriteString("----------------------------------------------------------------------------------\n")
	b.WriteString(fmt.Sprintf("Viewing routing table for node %v\n", local))
	for i, c := range state.RoutingTable.FindClosestContacts(local.ID, 20) {
		b.WriteString(fmt.Sprintf("%v: %v at distance %v\n", i, c, local.ID.CalcDistance(c.ID)))
	}
	b.WriteString("----------------------------------------------------------------------------------\n")
	b.WriteString(fmt.Sprintf("Last updated at %v", time.Now().Format("2006-01-02T15:04:05")))
	return b.String()
}

func main() {
	var myID *kademlia.KademliaID

	var origin = flag.Bool("origin", false, "should node be a bootstrap node")
	var listen = flag.String("listen", "0.0.0.0:4000", "which ip:port to listen for kademlia on")
	var s3listen = flag.String("s3listen", "0.0.0.0:8080", "which ip:port s3 should listen to")

	var bootstrapID = flag.String("bootstrap-id", "", "ID of node to bootstrap towards")
	var bootstrapIP = flag.String("bootstrap-ip", "", "IP of node to bootstrap towards")
	var consul = flag.Bool("consul", false, "bootstrap via consul")

	flag.Parse()

	if *origin {
		myID = kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	} else {
		rand.Seed(time.Now().UnixNano())
		myID = kademlia.NewRandomKademliaID()
	}

	if strings.Contains(*listen, "0.0.0.0") {
		parts := strings.Split(*listen, ":")
		log.Printf("[DEBUG] kadfs: Using local ip %v for assignment", GetLocalIP())
		*listen = GetLocalIP() + ":" + parts[1]
	}
	me := kademlia.NewContact(myID, *listen)
	myNetwork := kademlia.NewNetwork(&me)

	state := kademlia.NewKademliaState(me, myNetwork)

	// Starting all go routines
	go state.Network.Listen()
	go func() {
		for {
			// should probably be different go routines with different time for updates
			timer := time.NewTimer(150 * time.Second)
			<-timer.C
			log.Printf("[INFO] kadfs: Running republish, expire and replicate")
			go state.Replicate()
			go state.Republish()
			go state.Expire()
		}
	}()

	if *origin {
		log.Println("[INFO] kadfs: Running in origin mode, no bootstrap!")
	} else if *consul {
		log.Printf("[INFO] kadfs: Bootstrapping via consul")
		time.Sleep(2 * time.Second)

		client, err := api.NewClient(&api.Config{
			Address: "127.0.0.1:8500",
		})
		if err != nil {
			log.Panicf("[ERROR] kadfs: Unable to bootstrap via consul, error: %v", err)
		}
		parts := strings.Split(*listen, ":")
		address := parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Panicf("[ERROR] kadfs: Invalid port on listen for consul service registration")
		}
		tags := []string{fmt.Sprintf("kadfsid-%v", myID.String())}

		services, _, err := client.Catalog().Service("kadfs", "", &api.QueryOptions{})
		if err != nil {
			log.Panicf("[ERROR] kadfs: Unable to fetch services, error: %v", err)
		}
		if len(services) != 0 {
			rand.Seed(time.Now().Unix())
			service := services[rand.Intn(len(services))]

			consulBootstrapIP := fmt.Sprintf("%v:%v", service.ServiceAddress, service.ServicePort)
			consulBootstrapID := kademlia.NewKademliaID(strings.Replace(service.ServiceTags[0], "kadfsid-", "", 1))
			bootstrapNode := kademlia.NewContact(consulBootstrapID, consulBootstrapIP)

			state.Bootstrap(&bootstrapNode)
		} else {
			log.Printf("[INFO] kadfs: No services found in consul, registering and hoping someone will bootstrap towards me")
		}
		err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
			Name:    "kadfs",
			Address: address,
			Port:    port,
			Tags:    tags,
		})
		if err != nil {
			log.Panicf("[ERROR] kadfs: Failed to register service with consul, error: %v", err)
		} else {
			client.Agent().ServiceRegister(&api.AgentServiceRegistration{
				Name:    "kadfs-s3",
				Address: address,
				Port:    8080,
				Tags:    []string{"urlprefix-/"},
				Check: &api.AgentServiceCheck{
					TCP:      fmt.Sprintf("%v:%v", address, 8080),
					Interval: "10s",
					Timeout:  "2s",
				},
			})
			go func() {
				for {
					time.Sleep(10 * time.Second)
					client.KV().Put(&api.KVPair{
						Key:   fmt.Sprintf("routing-table/%v", *listen),
						Value: []byte(getRoutingTable(state)),
					}, nil)
				}
			}()
		}
	} else {
		log.Printf("[INFO] kadfs: Sleeping for 2 seconds to make sure the bootstrap node is up.")
		time.Sleep(2 * time.Second)

		raddr, err := net.ResolveUDPAddr("udp", *bootstrapIP)
		if err != nil {
			log.Fatalf("[ERROR] kadfs: Could not resolve %v", *bootstrapIP)
		}

		id2 := kademlia.NewKademliaID(*bootstrapID)
		bootstrapNode := kademlia.NewContact(id2, raddr.String()) // TODO: change

		// Should probably retry the boostrap a few times if we fail
		state.Bootstrap(&bootstrapNode)
	}

	// Listen for user input here whenever that gets implemented
	// fmt.Scanln()

	go s3.ConfigureAndListen(state, *s3listen)

	for {
		time.Sleep(15 * time.Second)
		examineRoutingTable(state)
	}

}
