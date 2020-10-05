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
	/*	--- OLD ----
	closestContacts := kademlia.network.routingTable.FindClosestContacts(target.ID, alpha)
	kademlia.PerformQuery(closestContacts, target)
		---- OLD ----*/
	
	// --- NEW ---
	
	var visitedList []string
	closestContacts := kademlia.network.routingTable.FindClosestContacts(target.ID, alpha)
	/*
	for _, c := range closestContacts{
		visitedList = append(visitedList, c.Address)
	}*/
	visitedList = append(visitedList, kademlia.network.routingTable.me.Address) //adds node itself to the visitedList in order to prevent it lookuping itself
	if(len(closestContacts)>0){ //prevent out of bounds on closest so far
		closestFromMe := closestContacts[0].ID.CalcDistance(target.ID)
		kademlia.PerformQuery(closestContacts, target, visitedList, closestFromMe)
	}
	// --- NEW ---

	// TODO (Node look up (Node Join))
	//	1. 	Async calls (Alpha decides how many?) to search for the contact in the 
	//		network (Using network.sendFindContactMessage).
	//	2. 	Check if the contact was found, return ip, UDP-port, node-id of the closest.
	//	3. 	Iterate with the response from step 2 if not found. 
}

func (kademlia *Kademlia) LookupData(hash string) {
	/*
		Should use LookupContact in order to find several nodes and check each of the nodes if they got any data "attached" to the given hash
	*/
}

func (kademlia *Kademlia) Store(data []byte) {
	/*
		Should just store the data on the node in a hashtable
	*/
}

func (kademlia *Kademlia) NodeJoin(address string) {
	contactToAdd := kademlia.network.SendNodeJoinMessage(address, kademlia.network.routingTable.me)
	kademlia.network.routingTable.AddContact(contactToAdd)
	kademlia.LookupContact(&kademlia.network.routingTable.me)
	//Add node lookup with this node as the target (Then the node should get an updated routing-table)
	
	// Comment back in to check for bucket elements
	/*
	for i := 0; i < 160 ; i++ {
		if(kademlia.network.routingTable.buckets[i].Len() > 0){
			fmt.Println(kademlia.network.routingTable.buckets[i].list.Front())
		}
		
	}
	*/
}
/*
	Performs the "iteration process" where the alpha-closest from the parent nodes get queried with their closest contacts
	For each of the alpha contacts, a thread starts and queries the contact for its alpha-closest. The result will be stored in a.
	The result-array (returnContacts) of each querie will be alpha*alpha in size and will then be used/manipulated to perform further queries.


	INITIATING NODE HAS FOUND K CONTACTS ALIVE
*/
func (kademlia *Kademlia) PerformQuery(contacts []Contact, target *Contact, visitedIPs []string, closestSoFar *KademliaID) {
	var returnContacts []Contact
	var a []Contact
	var queriedClosest *KademliaID
	
	/*	---- OLD ----
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
		}else{

		}
	}else{
		fmt.Println("Error: No contacts!")
	}	---- OLD ----*/
	
	//	---- NEW ----
	
	/*
		Start off by deleting the already visited nodes from the contact list (the list we will query from)
	*/
	for index, c := range contacts{
		_, found := Find(visitedIPs, c.Address)
		if found {
			contacts = DeleteByAddress(c.Address, contacts)	//WILL THIS BREAK? SINCE WE ARE LOOPING THROUGH THE SLICE ITSELF
		}
	}

	//PRINT ONLY!!
	for _, c := range contacts{
		fmt.Println("CONTACT TO REQUEST FROM: "+contacts.Address)
	}

	srtContact := kademlia.SortListBasedOnID(contacts, target)	//Needs to be sorted again after deletion. Stupid delete implementation.
	
	var count = 0 						//counter to prevent more than alpha concurrent calls
	for i := 0; i<len(srtContact); i++{	//loop on srtContact length in order to prevent out of bounds exception
		if count < alpha{
			fmt.Println("QUERIES THIS CONTACT FOR NODES: "+srtContact[i].String())
			a = <- kademlia.requestFromClosest(&srtContact[i], target)
			visitedIPs = append(visitedIPs, srtContact[i].Address)	//add the queried node's ip to the array of visited nodes ip's
			returnContacts = append(returnContacts,a...)
		}
		count = count+1	
	}

	//PRINT ONLY!!
	for _, ip := range visitedIPs{
		fmt.Println("THESE CONTACT IP'S HAVE BEEN VISITED: "+ip)
	}
	var contSize := 0
	for _, b := range kademlia.network.routingTable.buckets{
		contSize = contSize + b.Len()
	}
	fmt.Println("NUMBER OF CONTACTS IN ROUTING TABLE: "+ contSize)
	//PRINT ONLY!! ^^^

	for _, co := range returnContacts{
		kademlia.network.routingTable.AddContact(co)
	}
	sortedReturnContacts := kademlia.SortListBasedOnID(returnContacts, target)
	if(len(sortedReturnContacts)>0){ //prevent index out of bounds
		queriedClosest = sortedReturnContacts[0].ID.CalcDistance(target.ID)
		if(closestSoFar.Less(queriedClosest) && len(contacts) > 0){
			//Made no progress in regards of distance this iteration
			// WHAT SHOULD BE THE DIFFERENCE???????????????????????????
			kademlia.PerformQuery(sortedReturnContacts, target, visitedIPs, closestSoFar)
		}else{
			//Made progress in regards of distance this iteration
			kademlia.PerformQuery(sortedReturnContacts, target, visitedIPs, queriedClosest)
		}
	}else{

	}

	//	---- NEW ----
	
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

/*
	Finds the closest distance from the list of contacts with respect to the target
	Returns the index in the array of that contact, the contact itself
*/

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
	return 0, kademlia.network.routingTable.me, "Warning: Few contacts!"
	
}

/*
	Sorts a list based on the distance. Uses FindClosestDist n-times (size of list) and appends that result to the "result list"
	The closest-distance-contact get removed from the searched list for each time, hence the sort will work
*/

func (kademlia *Kademlia)SortListBasedOnID(contacts []Contact, target *Contact) []Contact{
	var newList []Contact
	for i := 0; i <= alpha; i++{	//Loops through index 0,1,2,3? Although alpha=3 implies that there should only be tree elements in the list
		index, contact, err := kademlia.FindClosestDist(contacts, target)
		if err != "" {
			
		}else{
			newList = append(newList, contact)
			contacts = DeleteFromContactList(contacts, index)
		}
		
	}
	return newList
}

/*
	Deletes an element at index i from a list
*/
func DeleteFromContactList(contacts []Contact, i int) []Contact{
	// Remove the element at index i from contacts.
	contacts[i] = contacts[len(contacts)-1] // Copy last element to index i.
	contacts = contacts[:len(contacts)-1]   // Truncate slice.
	return contacts
}

func DeleteByAddress(a string, contacts []Contact) []Contact{
	var con []Contact
	for i, c := range contacts{
		if c.Address == a{
			con = DeleteFromContactList(contacts, i)
			break
		}
	}
	return con
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}

// Creates a new kademlia struct
func InitKademlia(network *Network) *Kademlia{
	kademlia := &Kademlia{}
	kademlia.network = network
	return kademlia
}
