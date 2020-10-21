package d7024e

import (
	//"strconv"
	"fmt"
	"encoding/hex"
	"crypto/sha1"
	"regexp"
)	

type Kademlia struct {
	network *Network
}

const alpha = 3
const k = 20

/*
	1. 	Find alpha closest nodes (From routing table)
	2. 	Send request for their closest nodes to the target (From their routing tables)
	3. 	Repeat the second step on the nodes found until K nodes have been prompted,
		no more nodes to search or the node is found.
*/
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact{	
	fmt.Println("Looking up contact: " + target.ID.String())
	var visitedList []Contact
	closestContacts := kademlia.network.routingTable.FindClosestContacts(target.ID, alpha)
	visitedList = append(visitedList, kademlia.network.routingTable.me) //adds node itself to the visitedList in order to prevent it lookuping itself
	if(len(closestContacts)>0){ //prevent out of bounds on closest so far
		closestFromMe := closestContacts[0].ID.CalcDistance(target.ID)
		ret := kademlia.PerformQuery(closestContacts, target, visitedList, closestFromMe, 0)
		sortRet := kademlia.SortListBasedOnID(ret, target.ID)
		sortRet = DeleteByAddress(kademlia.network.routingTable.me.Address, sortRet)
		if(len(sortRet) > 0){
			fmt.Println("Closest contact found: " + sortRet[0].ID.String())
		}
		return sortRet
	}
	fmt.Println("No contacts found")
	return nil
}

func (kademlia *Kademlia) LookupData(hash string) string {
	/*
		Should use LookupContact in order to find several nodes and check each of the nodes if they got any data "attached" to the given hash
	
		1. 	Check that the hash is valid
		2. 	Check if the data is found on this node
		3.	Send request to look for the data on other nodes 
			1.	Lookup contact (key = target)
			2.	Send request to all contacts found in the lookup.
	*/
	var visitedList []Contact
	match, _ := regexp.MatchString(`[0-9a-fA-F]{40}`, hash)
	if !match {
		fmt.Println("Key not valid")
		return ""
	}

	value := kademlia.SearchKademliaStorage(hash)
	if value != nil {
		return string(value)
	}else{
		fmt.Println("####### I am searching for data in my network! ######")
		// Search for the value in kademlia network
		key := NewKademliaID(hash)
		visitedList = append(visitedList, kademlia.network.routingTable.me)
		closestContacts := kademlia.network.routingTable.FindClosestContacts(key, alpha) // Get the alpha closest contaccts to our key from routing table
		closestFromKey := closestContacts[0].ID.CalcDistance(key)
		data := kademlia.SearchByKey(closestContacts,key,visitedList,closestFromKey,key)
		if(data != nil){
			fmt.Println("####### I found data "+string(data)+" in my network! ######")
			return string(data)
		}else{
			return ""
		}
	}
}

func (kademlia *Kademlia) SearchKademliaStorage(hash string) []byte {
	key := NewKademliaID(hash)
	
	for _, e := range kademlia.network.hashtable {
		if e.Key.Equals(key) {
			return e.Data
		}
	}

	return nil
}

func (kademlia *Kademlia) Store(dataWritten []byte) {
	/*
		Should store the data on the node in a hashtable
		Should also store the data on the k-closest nodes to the hash (with respect to kademlia id)
	*/
	var a bool
	var result []bool
	h := sha1.New()
	h.Write(dataWritten)
	hexEncodedContent := hex.EncodeToString(h.Sum(nil))
	fmt.Println("This is the hash: "+hexEncodedContent)
	keyToAdd := NewKademliaID(hexEncodedContent)
	dataToAdd := &Data{
		Data: dataWritten,
		Key: keyToAdd,
	}
	//store internally
	kademlia.network.AddToHashTable(dataToAdd)
	//dummy contact
	dummyContact := NewContact(keyToAdd, "")
	//lookup for closest ID and send store-rpc
	closestContacts := kademlia.LookupContact(&dummyContact)
	for _,c := range closestContacts{
		fmt.Println("Storing at address: "+c.Address)
		a = <- kademlia.sendStoreToClosest(&c, dataToAdd)
		result = append(result, a)
	}
	for _,b := range result{
		if(b==false){
			fmt.Println("Failed at store")
		}
	}
}
/*
	1. Join the given node (Add it to routingtable and it adds this node to its routingtable)
	2. Run lookup on this node to get the closest neigbours to the routing table
*/
func (kademlia *Kademlia) NodeJoin(address string) {
	fmt.Println("Joining node with address: "+address)
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
	
	for  _, c := range contacts{
		if(probedContacts >= k){
			contacts = DeleteByAddress(c.Address, contacts)	
		}
	}

	srtContact := kademlia.SortListBasedOnID(contacts, target.ID)	//Needs to be sorted again after deletion. Bad delete implementation.
	
	var count = 0 						//counter to prevent more than alpha concurrent calls
	for i := 0; i<len(srtContact); i++{	//loop on srtContact length in order to prevent out of bounds exception
		if count < alpha{
			_, found := Find(visited, srtContact[i])
			if !found {	
				a = <- kademlia.requestFromClosest(&srtContact[i], target)
				visited = append(visited, srtContact[i])	//add the queried node's ip to the array of visited nodes ip's
				returnContacts = append(returnContacts,a...)
				probedContacts = probedContacts +1
			}else{
				continue
			}
			
		}else{
			break
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
	sortedReturnContacts := kademlia.SortListBasedOnID(returnContacts, target.ID)
	/*PRINT ONLY!!
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
	*/

	if(len(sortedReturnContacts)>0){ //prevent index out of bounds
		queriedClosest = sortedReturnContacts[0].ID.CalcDistance(target.ID)
		if(closestSoFar.Less(queriedClosest) && len(contacts) > 0){
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
	Search for the key in the network
	Closest contacts to targets are requested from remote nodes
	Send request to the nodes received recursively
*/
func (kademlia *Kademlia) SearchByKey(contacts []Contact, target *KademliaID, visited []Contact, closestSoFar *KademliaID, key *KademliaID) []byte{
	var returnContacts []Contact
	var a []Contact
	var data []byte
	var queriedClosest *KademliaID
	targetContact := NewContact(target,"")
	/*
		Start off by deleting the already visited nodes from the contact list (the list we will query from)
	*/
	for _, c := range contacts{
		_, found := Find(visited, c)
		if found {
			contacts = DeleteByAddress(c.Address, contacts)	
		}
	}
	
	//PRINT ONLY!!
	/*
	for _, c := range contacts{
		fmt.Println("CONTACT TO REQUEST FROM: "+c.Address)
	}
	*/

	srtContact := kademlia.SortListBasedOnID(contacts, target)	//Needs to be sorted again after deletion. Bad delete implementation.
	
	var count = 0 						//counter to prevent more than alpha concurrent calls
	for i := 0; i<len(srtContact); i++{	//loop on srtContact length in order to prevent out of bounds exception
		if count < alpha{
			a = <- kademlia.requestFromClosest(&srtContact[i], &targetContact)
			visited = append(visited, srtContact[i])	//add the queried node's ip to the array of visited nodes ip's
			returnContacts = append(returnContacts,a...)
			data = <- kademlia.requestDataFromClosest(&srtContact[i], key)
			if(data != nil){
				return data
			}
		}
		count = count+1	
	}
	for _, co := range returnContacts{
		if(kademlia.network.routingTable.me.ID.Equals(co.ID)){
			returnContacts = DeleteByAddress(co.Address, returnContacts)
		}else{
			kademlia.network.AddContHelper(co) // Add contacts to the routing table
		}
		
	}
	sortedReturnContacts := kademlia.SortListBasedOnID(returnContacts, target)


	if(len(sortedReturnContacts)>0){ //prevent index out of bounds
		queriedClosest = sortedReturnContacts[0].ID.CalcDistance(target)
		if(closestSoFar.Less(queriedClosest) && len(contacts) > 0){
			//Made no progress in regards of distance this iteration
			return kademlia.SearchByKey(sortedReturnContacts, target, visited, closestSoFar, key)
		}else{
			//Made progress in regards of distance this iteration
			return kademlia.SearchByKey(sortedReturnContacts, target, visited, queriedClosest, key)
		}
	}else{
		return nil
	}
}

/*
	Store the data on nodes closest to the key
*/
func (kademlia *Kademlia) sendStoreToClosest(contact *Contact, data *Data) <-chan bool{
	r := make(chan bool)
	go func(){
		defer close(r)
		reply:=kademlia.network.SendStoreMessage(contact.Address, data)
		r <- reply
	}()
	return r
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
*/
func (kademlia *Kademlia) requestDataFromClosest(contact *Contact, target *KademliaID) <-chan []byte {
	r := make(chan []byte)
	go func(){
		defer close(r)
		reply:=kademlia.network.SendFindDataMessage(contact.Address,target)
		r <- reply
	}()
	return r
}


/*
	Finds the closest distance from the list of contacts with respect to the target
	Returns the index in the array of that contact, the contact itself
*/

func (kademlia *Kademlia)FindClosestDist(contacts []Contact, target *KademliaID) (int,Contact, string){
	if len(contacts) > 0 {
		closest := contacts[0].ID
		contact := contacts[0]
		index := 0
		for i, c := range contacts{
			if c.ID.CalcDistance(target).Less(closest){
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

func (kademlia *Kademlia)SortListBasedOnID(contacts []Contact, target *KademliaID) []Contact{
	var newList []Contact
	var oldContacts []Contact
	oldContacts = contacts
	for i := 0; i < len(oldContacts); i++{	//Loops through index 0,1,2,3? Although alpha=3 implies that there should only be tree elements in the list
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
	con = contacts
	for i, c := range contacts{
		if (c.Address == a){
			con = DeleteFromContactList(contacts, i)
			return con
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
