package d7024e

type Kademlia struct {
	routingTable *RoutingTable
	network *Network
	// Routing table holds the contact information about this node 
	// It also has information about the bucket and holds information about contacts that this node knows are in the network.
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	//closestContacts := routingTable.FindClosestContacts(target.ID, 5) // 5 should be the size of the bucket size or alpha? 
	network.SendFindContactMessage(target)
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
func InitKademlia(me Contact, network *Network) *Kademlia{
	kademlia := &Kademlia{}
	kademlia.routingTable = NewRoutingTable(me)
	kademlia.network = network
	return kademlia
}
