package d7024e

<<<<<<< HEAD

=======
import (
	"fmt"
)
>>>>>>> 05d63cfa74dfa44783019d4c691272000ec6b4cd

type Kademlia struct {
	routingTable *RoutingTable
	network *Network
	// Routing table holds the contact information about this node 
	// It also has information about the bucket and holds information about contacts that this node knows are in the network.
}

<<<<<<< HEAD
const alpha = 3

func (kademlia *Kademlia) LookupContact(target *Contact) {
	closestContacts := kademlia.routingTable.FindClosestContacts(target.ID, alpha) // 3 should be the size of the bucket size or alpha? 
	kademlia.PerformQuery(closestContacts, target)
=======
func (kademlia *Kademlia) LookupContact(target *Contact) {
	closestContacts := kademlia.routingTable.FindClosestContacts(target.ID, 3) // 3 should be the size of the bucket size or alpha? 
	fmt.Println(len(closestContacts))
	for _, c:= range closestContacts{
		go func(){
			reqClosest := kademlia.network.SendFindContactMessage(&c, target)
			fmt.Println(reqClosest);
		}()
	}
>>>>>>> 05d63cfa74dfa44783019d4c691272000ec6b4cd
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

func (kademlia *Kademlia) PerformQuery(contacts []Contact, target *Contact) {
	var returnContacts []Contact
	var a []Contact
	for _, c := range contacts{
		a = <- kademlia.requestFromClosest(&c, target)
		returnContacts=append(returnContacts,a...)
	}
	for _, co := range returnContacts{
		kademlia.routingTable.AddContact(co)
	}
	sortedReturnContacts := SortListBasedOnID(returnContacts, target)
	if(contacts[0].ID.Less(sortedReturnContacts[0].ID)){
		kademlia.PerformQuery(sortedReturnContacts, target)
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

func FindClosestDist(contacts []Contact, target *Contact) (int,Contact){
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
	return index,contact
}

func SortListBasedOnID(contacts []Contact, target *Contact) []Contact{
	var newList []Contact
	for i := 0; i <= alpha; i++{
		index, contact := FindClosestDist(contacts, target)
		newList[i] = contact
		contacts = DeleteFromContactList(contacts, index)
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
func InitKademlia(rt *RoutingTable, network *Network) *Kademlia{
	kademlia := &Kademlia{}
	kademlia.routingTable = rt
	kademlia.network = network
	return kademlia
}
