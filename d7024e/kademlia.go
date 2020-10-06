package d7024e

import (
	"strconv"
	"fmt"
	"encoding/hex"
)	

type Kademlia struct {
	network *network
	hashtable *[]Data
}

type Data struct {
	data []byte
	key *KademliaID
}

const alpha = 3
<<<<<<< HEAD
const bucketSize = 20

=======
>>>>>>> ec607a6a15ec94ece86b401a7f83032d26f65f76
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact{	
	var visitedList []Contact
	closestContacts := kademlia.network.routingTable.FindClosestContacts(target.ID, alpha)
	visitedList = append(visitedList, kademlia.network.routingTable.me) //adds node itself to the visitedList in order to prevent it lookuping itself
	if(len(closestContacts)>0){ //prevent out of bounds on closest so far
		closestFromMe := closestContacts[0].ID.CalcDistance(target.ID)
		ret := kademlia.PerformQuery(closestContacts, target, visitedList, closestFromMe, 0)
		sortRet := kademlia.SortListBasedOnID(ret, target)
		sortRet = DeleteByAddress(kademlia.network.routingTable.me.Address, sortRet)
		fmt.Println("################# Sorted returned list ####################")
		for _, c := range sortRet{
			fmt.Println(c.ID.String())
		}
		return sortRet
	}
	return nil
}

func (kademlia *Kademlia) LookupData(hash string) {
	/*
		Should use LookupContact in order to find several nodes and check each of the nodes if they got any data "attached" to the given hash
	*/
}

func (kademlia *Kademlia) Store(data []byte) {
	/*
		Should store the data on the node in a hashtable
		Should also store the data on the k-closest nodes to the hash (with respect to kademlia id)
	*/
	h := sha1.New()
	h.Write(data)
	hexEncodedContent = hex.EncodeToString(h.Sum(nil))
	keyToAdd := NewKademliaID(hexEncodedContent)
	dataToAdd := &Data{
		data h.Sum(nil)
		key keyToAdd
	}
	//store internally
	kademlia.hashtable = append(kademlia.hashtable, dataToAdd)
	//dummy contact
	dummyContact := NewContact(keyToAdd, "0.0.0.0")
	//lookup for closest ID
	closestContacts := kademlia.network.routingTable.FindClosestContacts(target.ID, alpha)

}

func (kademlia *Kademlia) NodeJoin(address string) {
	contactToAdd := kademlia.network.SendNodeJoinMessage(address, kademlia.network.routingTable.me)
	kademlia.network.AddContHelper(contactToAdd)
	kademlia.LookupContact(&kademlia.network.routingTable.me)
}

/*
	Performs the "iteration process" where the alpha-closest from the parent nodes get queried with their closest contacts
	For each of the alpha contacts, a thread starts and queries the contact for its alpha-closest. The result will be stored in a.
	The result-array (returnContacts) of each querie will be alpha*alpha in size and will then be used/manipulated to perform further queries.


	INITIATING NODE HAS FOUND K CONTACTS ALIVE
*/
func (kademlia *Kademlia) PerformQuery(contacts []Contact, target *Contact, visited []Contact, closestSoFar *KademliaID, probedContacts int) []Contact{
	var returnContacts []Contact
	var a []Contact
	var queriedClosest *KademliaID
	/*
		Start off by deleting the already visited nodes from the contact list (the list we will query from)
	*/
	for _, c := range contacts{
		_, found := Find(visited, c)
		if found {
			contacts = DeleteByAddress(c.Address, contacts)	//WILL THIS BREAK? SINCE WE ARE LOOPING THROUGH THE SLICE ITSELF
		}
	}
	fmt.Println("PROBED CONTACTS: "+strconv.Itoa(probedContacts))
	for  _, c := range contacts{
		if(probedContacts >= bucketSize){
			contacts = DeleteByAddress(c.Address, contacts)	//WILL THIS BREAK? SINCE WE ARE LOOPING THROUGH THE SLICE ITSELF
			return visited
		}
	}
	//PRINT ONLY!!
	for _, c := range contacts{
		fmt.Println("CONTACT TO REQUEST FROM: "+c.Address)
	}

	srtContact := kademlia.SortListBasedOnID(contacts, target)	//Needs to be sorted again after deletion. Stupid delete implementation.
	
	var count = 0 						//counter to prevent more than alpha concurrent calls
	for i := 0; i<len(srtContact); i++{	//loop on srtContact length in order to prevent out of bounds exception
		if count < alpha{
			fmt.Println("QUERIES THIS CONTACT FOR NODES: "+srtContact[i].String())
			a = <- kademlia.requestFromClosest(&srtContact[i], target)
			visited = append(visited, srtContact[i])	//add the queried node's ip to the array of visited nodes ip's
			returnContacts = append(returnContacts,a...)
			probedContacts = probedContacts +1
		}
		count = count+1	
	}
	for _, co := range returnContacts{
		if(kademlia.network.routingTable.me.ID.Equals(co.ID)){
			returnContacts = DeleteByAddress(co.Address, returnContacts)
		}else{
			kademlia.network.AddContHelper(co)
		}
		
	}
	sortedReturnContacts := kademlia.SortListBasedOnID(returnContacts, target)
	//PRINT ONLY!!
	for _, c := range sortedReturnContacts{
		fmt.Println("Returned contact: "+c.Address)
	}
	for _, con := range visited{
		fmt.Println("THESE CONTACT IP'S HAVE BEEN VISITED: "+con.Address)
	}
	var contSize = 0
	for _, b := range kademlia.network.routingTable.buckets{
		contSize = contSize + b.Len()
	}
	fmt.Println("NUMBER OF CONTACTS IN ROUTING TABLE: "+ strconv.Itoa(contSize))
	//PRINT ONLY!! ^^^


	if(len(sortedReturnContacts)>0){ //prevent index out of bounds
		queriedClosest = sortedReturnContacts[0].ID.CalcDistance(target.ID)
		if(closestSoFar.Less(queriedClosest) && len(contacts) > 0){
			//Made no progress in regards of distance this iteration
			// WHAT SHOULD BE THE DIFFERENCE???????????????????????????
			
			return kademlia.PerformQuery(sortedReturnContacts, target, visited, closestSoFar, probedContacts)
		}else{
			//Made progress in regards of distance this iteration
			return kademlia.PerformQuery(sortedReturnContacts, target, visited, queriedClosest, probedContacts)
		}
	}else{
		return visited
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
func Find(slice []Contact, val Contact) (int, bool) {
    for i, item := range slice {
        if item.Address == val.Address {
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
