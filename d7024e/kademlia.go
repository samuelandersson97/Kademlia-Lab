package d7024e

type Kademlia struct {
	// Kademlia id
	// Bucketd
	routingTable *RoutingTable
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	closestContacts := routingTable.FindClosestContacts(target.ID, 5) // 5 should be the size of the bucket or alpha? 
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func InitKademlia(routingtable *RoutingTable) *Kademlia{
	kademlia := &Kademlia{}
	kademlia.contact = routingTable
	return kademlia
}