package d7024e

import (
	"fmt"
)

type Kademlia struct {
	routingTable *RoutingTable
	network *Network
	// Routing table holds the contact information about this node 
	// It also has information about the bucket and holds information about contacts that this node knows are in the network.
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	closestContacts := kademlia.routingTable.FindClosestContacts(target.ID, 1) // 3 should be the size of the bucket size or alpha? 
	fmt.Println(len(closestContacts))
	for i := 0; i<len(closestContacts); i++{
		
		go kademlia.network.SendFindContactMessage(&closestContacts[i], target)
		fmt.Println("reqClosest");
	}
	// TODO (Node look up (Node Join))
	//	1. 	Async calls (Alpha decides how many?) to search for the contact in the 
	//		network (Using network.sendFindContactMessage).
	//	2. 	Check if the contact was found, return ip, UDP-port, node-id of the closest.
	//	3. 	Iterate with the response from step 2 if not found. 
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

// Creates a new kademlia struct
func InitKademlia(rt *RoutingTable, network *Network) *Kademlia{
	kademlia := &Kademlia{}
	kademlia.routingTable = rt
	kademlia.network = network
	return kademlia
}
