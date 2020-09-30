package d7024e

import (
	"fmt"
)
	

type Kademlia struct {
	network *Network
	// Routing table holds the contact information about this node 
	// It also has information about the bucket and holds information about contacts that this node knows are in the network.
}

const alpha = 3

func (kademlia *Kademlia) LookupContact(target *Contact) {
	closestContacts := kademlia.network.routingTable.FindClosestContacts(target.ID, alpha) // 3 should be the size of the bucket size or alpha? 
	kademlia.PerformQuery(closestContacts, target)
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

func (kademlia *Kademlia) NodeJoin(address string) {
	contactToAdd := kademlia.network.SendNodeJoinMessage(address, kademlia.network.routingTable.me)
	kademlia.network.routingTable.AddContact(contactToAdd)
	//Add node lookup with this node as the target (Then the node should get an updated routing-table)
	
	/* // Comment back in to check for bucket elements
	for i := 0; i < 160 ; i++ {
		if(kademlia.network.routingTable.buckets[i].Len() > 0){
			fmt.Println(kademlia.network.routingTable.buckets[i].list.Front())
		}
		
	}*/
	
}

func (kademlia *Kademlia) PerformQuery(contacts []Contact, target *Contact) {
	var returnContacts []Contact
	var a []Contact
	for _, c := range contacts{
		a = <- kademlia.requestFromClosest(&c, target)
		returnContacts=append(returnContacts,a...)
	}
	for _, co := range returnContacts{
		kademlia.network.routingTable.AddContact(co)
	}
	sortedReturnContacts := kademlia.SortListBasedOnID(returnContacts, target)
	if(len(sortedReturnContacts) > 0 && len(contacts) > 0){
		if(contacts[0].ID.Less(sortedReturnContacts[0].ID)){
			kademlia.PerformQuery(sortedReturnContacts, target)
		}
	}else{
		fmt.Println("Error: No contacts!")
	}
	
}

/*
	Requests the closest alpha-contacts from a given contact.
	Returns a list of contacts.

*/
func (kademlia *Kademlia) requestFromClosest(contact *Contact, target *Contact) <-chan []Contact {
	r := make(chan []Contact)
	go func(){
		defer close(r)
		reply:=kademlia.network.SendFindContactMessage(contact,target)
		r <- reply
	}()
	return r
}

func (kademlia *Kademlia)FindClosestDist(contacts []Contact, target *Contact) (int,Contact, string){
	if len(contacts) > 0 {
		closest := contacts[0].ID
		contact := contacts[0]
		index := 0
		for i, c := range contacts{
			if c.ID.CalcDistance(target.ID).Less(closest){
				closest = c.ID
				contact = c
				index = i
			}
		}
		return index,contact,""
	}
	return 0, kademlia.network.routingTable.me, "Error: No contacts!"
	
}

func (kademlia *Kademlia)SortListBasedOnID(contacts []Contact, target *Contact) []Contact{
	var newList []Contact
	for i := 0; i <= alpha; i++{
		index, contact, err := kademlia.FindClosestDist(contacts, target)
		if err != "" {
			fmt.Println(err)
		}else{
			newList = append(newList, contact)
			contacts = DeleteFromContactList(contacts, index)
		}
		
	}
	return newList
}

func DeleteFromContactList(contacts []Contact, i int) []Contact{
	// Remove the element at index i from contacts.
	contacts[i] = contacts[len(contacts)-1] // Copy last element to index i.
	contacts = contacts[:len(contacts)-1]   // Truncate slice.
	return contacts
}

// Creates a new kademlia struct
func InitKademlia(network *Network) *Kademlia{
	kademlia := &Kademlia{}
	kademlia.network = network
	return kademlia
}
