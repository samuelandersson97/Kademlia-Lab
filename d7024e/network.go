package d7024e

import(
	"net"
	"fmt"
	"time"
	"encoding/json"
)

type Network struct {
	routingTable *RoutingTable
	hashtable []*Data
}

type Protocol struct {
	Rpc string			//rpc-type (node join, node lookup, ping etc...)
	Contacts []Contact  //only used for return from node lookup
	Hash string			//used for get
	Data []byte			//sending data needed for given rpc
	Message string		//used for sent/recieve of the rpc
}

type Data struct {
	Data []byte
	Key *KademliaID
}

type LookupData struct {
	Target *Contact
	Sender Contact
}

const T_OUT = time.Millisecond*1000.

/*
	Extract the sent message and create different fucntions that handles different types of messages (PING, FIND_NODE, etc...)
*/
func (network *Network) Listen(ip string) {
	//Returns an address of the UDP end point. 'udp4' indicates that only IPv4-addresses are being resolved
	udpEndPoint, err := net.ResolveUDPAddr("udp4",ip)
	if err != nil {
		fmt.Println(err)
	}
	//Listens for packets on the (ONLY!!) LOCAL network. 'udp4' indicates that only IPv4-addresses are taken into account when it comes to listening for packets, returns a connection
	c, err := net.ListenUDP("udp4", udpEndPoint)
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()
	//creates buffer with maximum length of 8192
	messageBuffer := make([]byte, 8192)
	responseProtocol := Protocol{}
	for{
		//Adds the message from the UDP-channel in the message-buffer. Returns the size of the message and the adress of the sender
		size, senderAddress, err := c.ReadFromUDP(messageBuffer)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(messageBuffer[:size], &responseProtocol)
		_ = network.DecodeRPC(&responseProtocol, senderAddress, c)
	}
}

func (network *Network) SendPingMessage(contact *Contact) bool {
	//Returns an address of the UDP end point. 'udp4' indicates that only IPv4-addresses are being resolved
	udpEndPoint, err := net.ResolveUDPAddr("udp4",contact.Address)
	fmt.Println("Trying to ping: "+contact.Address)
	if err != nil {
		fmt.Println(err)
	}
	// Starts up a UDP-connection to the resolved UDP-address 
	c, err := net.DialUDP("udp4",nil, udpEndPoint)
	if err != nil {
		fmt.Println(err)
	}
	meContact := ContactToByteArray(&network.routingTable.me)
	pingMessage := CreateProtocol("PING", nil, "", meContact, "PING_SENT")
	c.SetDeadline(time.Now().Add(T_OUT))
	defer c.Close()
	_, e := c.Write(pingMessage)
	if e != nil {
		fmt.Println(err)
	}
	responseBuffer := make([]byte, 8192)
	size, senderAddress, err := c.ReadFromUDP(responseBuffer)
	receivedPing := Protocol{}
	if err != nil {
		fmt.Println("TIMEOUT")
		fmt.Println("Failed to ping: "+contact.Address)
		return false
	}else{
		json.Unmarshal(responseBuffer[:size], &receivedPing)
		_ = network.DecodeRPC(&receivedPing, senderAddress, c)
		fmt.Println("Succesfully pinged: "+contact.Address)
		return true
	}


	
}

func (network *Network) SendFindContactMessage(contact *Contact, target *Contact) []Contact {
	byteArr := LookupDataToByteArray(target, network.routingTable.me)
	lookupData := LookupData{}
	json.Unmarshal(byteArr[:len(byteArr)], &lookupData)
	//Returns an address of the UDP end point. 'udp4' indicates that only IPv4-addresses are being resolved
	udpEndPoint, err := net.ResolveUDPAddr("udp4",contact.Address)
	if err != nil {
		fmt.Println(err)
	}
	// Starts up a UDP-connection to the resolved UDP-address 
	c, err := net.DialUDP("udp4",nil, udpEndPoint)
	if err != nil {
		fmt.Println(err)
	}
	lookupMessage := CreateProtocol("NODE_LOOKUP", nil, "", byteArr, "NODE_LOOKUP_SENT")
	c.SetDeadline(time.Now().Add(T_OUT))
	defer c.Close()
	_, e := c.Write(lookupMessage)
	if e != nil {
		fmt.Println(err)
	}
	responseBuffer := make([]byte, 8192)
	size, senderAddress, err := c.ReadFromUDP(responseBuffer)
	receivedLookup := Protocol{}
	if err != nil {
		fmt.Println(err)
		return nil
	}else{
		json.Unmarshal(responseBuffer[:size], &receivedLookup)
		contactProt := network.DecodeRPC(&receivedLookup, senderAddress, c)
		return contactProt.Contacts
	}
}

func (network *Network) SendNodeJoinMessage(address string, me Contact) Contact {
	joinData := ContactToByteArray(&me)
	udpEndPoint, err := net.ResolveUDPAddr("udp4",address)
	if err != nil {
		fmt.Println(err)
	}
	// Starts up a UDP-connection to the resolved UDP-address 
	c, err := net.DialUDP("udp4",nil, udpEndPoint)
	if err != nil {
		fmt.Println(err)
	}
	joinMessage := CreateProtocol("NODE_JOIN", nil, "", joinData, "NODE_JOIN_SENT")
	c.SetDeadline(time.Now().Add(T_OUT))
	defer c.Close()
	_, e := c.Write(joinMessage)
	if e != nil {
		fmt.Println(err)
	}
	responseBuffer := make([]byte, 8192)
	size, senderAddress, err := c.ReadFromUDP(responseBuffer)
	receivedJoin := Protocol{}
	contactInformation := Contact{}
	if err != nil {
		fmt.Println(err)
		return contactInformation
	}
	json.Unmarshal(responseBuffer[:size], &receivedJoin)
	responseProtocol := network.DecodeRPC(&receivedJoin, senderAddress, c)
	json.Unmarshal(responseProtocol.Data[:len(responseProtocol.Data)], &contactInformation)
	return contactInformation
	
}

func (network *Network) SendFindDataMessage(address string, key *KademliaID) []byte{
	udpEndPoint, err := net.ResolveUDPAddr("udp4",address)
	if err != nil {
		fmt.Println(err)
	}
	// Starts up a UDP-connection to the resolved UDP-address 
	c, err := net.DialUDP("udp4",nil, udpEndPoint)
	if err != nil {
		fmt.Println(err)
	}
	dataToSend, _ := json.Marshal(key.String())
	findDataMessage := CreateProtocol("NODE_VALUE", nil, "", dataToSend, "NODE_VALUE_SENT")
	c.SetDeadline(time.Now().Add(T_OUT))
	defer c.Close()
	_, e := c.Write(findDataMessage)
	if e != nil {
		fmt.Println(err)
	}
	responseBuffer := make([]byte, 8192)
	size, senderAddress, err := c.ReadFromUDP(responseBuffer)
	findDataResponse := Protocol{}
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(responseBuffer[:size], &findDataResponse)
	responseProtocol := network.DecodeRPC(&findDataResponse, senderAddress, c)
	if(responseProtocol.Message=="NODE_VALUE_SUCCESS"){
		return responseProtocol.Data
	}else{
		return nil
	}
	return nil
}

func (network *Network) SendStoreMessage(address string, data *Data) bool{
	dataToSend := DataToByteArray(data)
	udpEndPoint, err := net.ResolveUDPAddr("udp4",address)
	if err != nil {
		fmt.Println(err)
	}
	// Starts up a UDP-connection to the resolved UDP-address 
	c, err := net.DialUDP("udp4",nil, udpEndPoint)
	if err != nil {
		fmt.Println(err)
	}
	storeMessage := CreateProtocol("NODE_STORE", nil, "", dataToSend, "NODE_STORE_SENT")
	c.SetDeadline(time.Now().Add(T_OUT))
	defer c.Close()
	_, e := c.Write(storeMessage)
	if e != nil {
		fmt.Println(err)
	}
	responseBuffer := make([]byte, 8192)
	size, senderAddress, err := c.ReadFromUDP(responseBuffer)
	storeResponse := Protocol{}
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(responseBuffer[:size], &storeResponse)
	responseProtocol := network.DecodeRPC(&storeResponse, senderAddress, c)
	if(responseProtocol.Message=="NODE_STORE_ACCEPTED"){
		return true
	}else{
		return false
	}
}

func (network *Network) DecodeRPC(prot *Protocol, senderAddress *net.UDPAddr, connection *net.UDPConn) *Protocol{
	if(prot.Rpc == "PING"){
		return network.PingHandler(prot, senderAddress, connection)
	}else if(prot.Rpc == "NODE_LOOKUP"){
		return network.LookupHandler(prot, senderAddress, connection)
	}else if(prot.Rpc == "NODE_JOIN"){
		return network.JoinHandler(prot, senderAddress, connection)
	}else if(prot.Rpc == "NODE_VALUE"){
		return network.FindDataHandler(prot, senderAddress, connection)
	}else if(prot.Rpc == "NODE_STORE"){
		return network.StoreHandler(prot, senderAddress, connection)
	}else{
		fmt.Println("ERROR. RPC TYPE COULD NOT BE FOUND")
		return nil
	}
}
func (network *Network) FindDataHandler(prot *Protocol, responseAddr *net.UDPAddr, connection *net.UDPConn) *Protocol{
	if(prot.Message == "NODE_VALUE_SENT"){
		key := ""
		var response []byte
		var protToSend []byte
		json.Unmarshal(prot.Data[:len(prot.Data)], &key)
		fmt.Println("Searching for data with key: "+ key)
		for _, e := range network.hashtable {
			if(e.Key.String() == key) {
				fmt.Println("###### Data found: " + string(e.Data) + " ######")
				response = e.Data
				break
			}
		}

		if(response != nil){
			protToSend = CreateProtocol("NODE_VALUE", nil, "", response, "NODE_VALUE_SUCCESS")
		}else{
			protToSend = CreateProtocol("NODE_VALUE", nil, "", nil, "NODE_VALUE_REJECTED")
		}
		_, e := connection.WriteToUDP(protToSend, responseAddr)
		if e != nil{
			fmt.Println(e)
		}
		return prot
	}else{
		return prot
	}
}

func (network *Network) StoreHandler(prot *Protocol, responseAddr *net.UDPAddr, connection *net.UDPConn) *Protocol{
	if(prot.Message == "NODE_STORE_SENT"){
		sentData := Data{}
		var protToSend []byte
		json.Unmarshal(prot.Data[:len(prot.Data)], &sentData)
		success := network.AddToHashTable(&sentData)
		if(success){
			for _, d := range network.hashtable {
				fmt.Println(string(d.Data))
			}
			protToSend = CreateProtocol("NODE_STORE", nil, "", nil, "NODE_STORE_ACCEPTED")
		}else{
			protToSend = CreateProtocol("NODE_STORE", nil, "", nil, "NODE_STORE_REJECTED")
		}
		_, e := connection.WriteToUDP(protToSend, responseAddr)
		if e != nil{
			fmt.Println(e)
		}
		return prot
	}else{
		return prot
	}
}

func (network *Network) JoinHandler(prot *Protocol, responseAddr *net.UDPAddr, connection *net.UDPConn) *Protocol{
	if(prot.Message == "NODE_JOIN_SENT"){
		sendContact := Contact{}
		json.Unmarshal(prot.Data[:len(prot.Data)], &sendContact)
		
		fmt.Println("Node with id "+ sendContact.ID.String() +" joined this node.")
		
		//Respond with my own contact
		meContact := ContactToByteArray(&network.routingTable.me)
		joinProtocolResponse := CreateProtocol("NODE_JOIN",nil,"",meContact,"NODE_JOIN_RESPONSE")
		_, e := connection.WriteToUDP(joinProtocolResponse, responseAddr)
		if e != nil{
			fmt.Println(e)
		}
		network.AddContHelper(sendContact) 
		return prot
	}else if(prot.Message == "NODE_JOIN_RESPONSE"){
		return prot
	}
	return nil
}

func (network *Network) LookupHandler(prot *Protocol, responseAddr *net.UDPAddr, connection *net.UDPConn) *Protocol{
	if(prot.Message == "NODE_LOOKUP_SENT"){
		lookupData := LookupData{}
		json.Unmarshal(prot.Data[:len(prot.Data)], &lookupData)
		closestContactsArray := network.routingTable.FindClosestContacts(lookupData.Target.ID, 3)
		/*for _,c := range closestContactsArray{
			fmt.Println("This contact was one of the closest: "+c.Address)
		}*/
		lookupProtocolResponse := CreateProtocol("NODE_LOOKUP",closestContactsArray,"",prot.Data,"NODE_LOOKUP_RESPONSE")
		_, e := connection.WriteToUDP(lookupProtocolResponse, responseAddr)
		if e != nil{
			fmt.Println(e)
		}
		network.AddContHelper(lookupData.Sender)
		return prot
	}else if(prot.Message == "NODE_LOOKUP_RESPONSE"){
		return prot
	}
	return nil
}

func (network *Network) PingHandler(prot *Protocol, responseAddr *net.UDPAddr, connection *net.UDPConn) *Protocol{
	/*
		What if a node is dead? The ping should be able to timeout somehow?
		There is "SetDeadline"-stuff in the documentation for net, check it out
	*/
	if(prot.Message == "PING_SENT"){
		contactToAdd := Contact{}
		json.Unmarshal(prot.Data[:len(prot.Data)], &contactToAdd)
		pingResponseRPC := CreateProtocol("PING",nil,"",nil,"PING_RESPONSE")
		_, e := connection.WriteToUDP(pingResponseRPC, responseAddr)
		if e != nil {
			fmt.Println(e)
		}
		network.AddContHelper(contactToAdd)
		return prot
	}else if(prot.Message == "PING_RESPONSE"){
		return prot
	}
	return nil
}

func ContactToByteArray(contact *Contact) []byte {
	contactByteArray, err := json.Marshal(contact)
	if err != nil{
		fmt.Println(err)
		return nil
	}
	return contactByteArray
}

func DataToByteArray(data *Data) []byte {
	dataByteArray, err := json.Marshal(data)
	if err != nil{
		fmt.Println(err)
		return nil
	}
	return dataByteArray
}

func LookupDataToByteArray(targetCont *Contact, senderCont Contact) []byte {
	data := &LookupData{
		Target: targetCont,
		Sender: senderCont}
	dataByteArray, err := json.Marshal(data)
	if err != nil{
		return nil
	}
	return dataByteArray
}

func CreateProtocol(rpcToSend string, contactsArr []Contact, hashToSend string, dataToSend []byte, messageToSend string) []byte{
	protocol := &Protocol{
		Rpc: rpcToSend,
		Contacts: contactsArr,
		Hash: hashToSend,
		Data: dataToSend,
		Message: messageToSend}
	prot, err := json.Marshal(protocol)
	if err != nil{
		fmt.Println(err)
		return nil
	}
	return prot
}

func InitNetwork(rt *RoutingTable) *Network{
	network := &Network{}
	network.routingTable = rt
	return network
}

//Gets the nodes local IP by dialing Google's server and returns the IP as a string
func GetOutboundIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        fmt.Println(err)
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP.String()
}


// Add to routing table 
func (network *Network) AddContHelper(contact Contact){
	if(contact.ID != nil){
		if(!network.routingTable.CheckIfFull(contact.ID)){
			network.routingTable.AddContact(contact) 
		}else{
			bucket := network.routingTable.GetBucket(contact.ID)
			cont := Contact{}
			cont = bucket.list.Back().Value.(Contact)
			if(network.SendPingMessage(&cont)){
				network.routingTable.AddContact(bucket.list.Back().Value.(Contact)) //see AddContact in bucket
			}else{
				network.routingTable.AddContact(contact)
			}
		}
	}
}

func (network *Network) AddToHashTable(data *Data) bool{
	_, b := FindInHash(network.hashtable, data)
	if(!b){
		network.hashtable = append(network.hashtable, data)
		return true
	}
	return false
}

//checks only if the value exists, not the key
func FindInHash(slice []*Data, val *Data) (int, bool) {
    for i, item := range slice {
        if(string(item.Data) == string(val.Data)){
            return i, true
        }
    }
    return -1, false
}
