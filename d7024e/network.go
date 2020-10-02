package d7024e

import(
	"net"
	"strconv"
	"fmt"
	"encoding/json"
)

/*
	Here we should have a pointer to a Kademlia "object". The Kademlia object itself contains a RoutingTable.
*/

type Network struct {
	routingTable *RoutingTable
}

type Protocol struct {
	Rpc string
	Contacts []Contact
	Hash string
	Data []byte
	Message string
}

/*
	Extract the sent message and create different fucntions that handles different types of messages (PING, FIND_NODE, etc...)
*/
func (network *Network) Listen(ip string, port int) {
	adrPort := ip+":"+strconv.Itoa(port)
	fmt.Println("Listening at "+adrPort+ ".....")
	//Returns an address of the UDP end point. 'udp4' indicates that only IPv4-addresses are being resolved
	udpEndPoint, err := net.ResolveUDPAddr("udp4",adrPort)
	if err != nil {
		fmt.Println("LISTEN ERROR: 1")
		fmt.Println(err)
	}
	//Listens for packets on the (ONLY!!) LOCAL network. 'udp4' indicates that only IPv4-addresses are taken into account when it comes to listening for packets, returns a connection
	c, err := net.ListenUDP("udp4", udpEndPoint)
	if err != nil {
		fmt.Println("LISTEN ERROR: 2")
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
			fmt.Println("LISTEN ERROR: 3")
			fmt.Println(err)
		}
		json.Unmarshal(messageBuffer[:size], &responseProtocol)
		_ = network.DecodeRPC(&responseProtocol, senderAddress, c)
	}
}

/*
	Clean up the trace prints, send useful messages instead of this (should be wrapped in our struct mentioned above)
*/

func (network *Network) SendPingMessage(contact *Contact) {
	//Returns an address of the UDP end point. 'udp4' indicates that only IPv4-addresses are being resolved
	udpEndPoint, err := net.ResolveUDPAddr("udp4",contact.Address+":1111")
	if err != nil {
		fmt.Println("SEND ERROR: 1")
		fmt.Println(err)
	}
	// Starts up a UDP-connection to the resolved UDP-address 
	c, err := net.DialUDP("udp4",nil, udpEndPoint)
	if err != nil {
		fmt.Println("SEND ERROR: 2")
		fmt.Println(err)
	}
	pingMessage := CreateProtocol("PING", nil, "", nil, "PING_SENT")
	defer c.Close()
	_, e := c.Write(pingMessage)
	if e != nil {
		fmt.Println("SEND ERROR: 3")
		fmt.Println(err)
	}
	responseBuffer := make([]byte, 8192)
	size, senderAddress, err := c.ReadFromUDP(responseBuffer)
	receivedPing := Protocol{}
	if err != nil {
		fmt.Println("SEND ERROR: 4")
		fmt.Println(err)
	}else{
		json.Unmarshal(responseBuffer[:size], &receivedPing)
		_ = network.DecodeRPC(&receivedPing, senderAddress, c)
	}
	
}

func (network *Network) SendFindContactMessage(contact *Contact, target *Contact) []Contact {
	targetArr := ContactToByteArray(target)
	//Returns an address of the UDP end point. 'udp4' indicates that only IPv4-addresses are being resolved
	udpEndPoint, err := net.ResolveUDPAddr("udp4",contact.Address+":1111")
	if err != nil {
		fmt.Println("SEND ERROR: 1")
		fmt.Println(err)
	}
	// Starts up a UDP-connection to the resolved UDP-address 
	c, err := net.DialUDP("udp4",nil, udpEndPoint)
	if err != nil {
		fmt.Println("SEND ERROR: 2")
		fmt.Println(err)
	}
	lookupMessage := CreateProtocol("NODE_LOOKUP", nil, "", targetArr, "NODE_LOOKUP_SENT")
	defer c.Close()
	_, e := c.Write(lookupMessage)
	if e != nil {
		fmt.Println("SEND ERROR: 3")
		fmt.Println(err)
	}
	responseBuffer := make([]byte, 8192)
	size, senderAddress, err := c.ReadFromUDP(responseBuffer)
	receivedLookup := Protocol{}
	if err != nil {
		fmt.Println("SEND ERROR: 4")
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
	udpEndPoint, err := net.ResolveUDPAddr("udp4",address+":1111")
	if err != nil {
		fmt.Println("SEND ERROR: 1")
		fmt.Println(err)
	}
	// Starts up a UDP-connection to the resolved UDP-address 
	c, err := net.DialUDP("udp4",nil, udpEndPoint)
	if err != nil {
		fmt.Println("SEND ERROR: 2")
		fmt.Println(err)
	}
	joinMessage := CreateProtocol("NODE_JOIN", nil, "", joinData, "NODE_JOIN_SENT")
	defer c.Close()
	_, e := c.Write(joinMessage)
	if e != nil {
		fmt.Println("SEND ERROR: 3")
		fmt.Println(err)
	}
	responseBuffer := make([]byte, 8192)
	size, senderAddress, err := c.ReadFromUDP(responseBuffer)
	receivedJoin := Protocol{}
	if err != nil {
		fmt.Println("SEND ERROR: 4")
		fmt.Println(err)
		
	}
	json.Unmarshal(responseBuffer[:size], &receivedJoin)
	responseProtocol := network.DecodeRPC(&receivedJoin, senderAddress, c)
	contactInformation := Contact{}
	json.Unmarshal(responseProtocol.Data[:len(responseProtocol.Data)], &contactInformation)
	return contactInformation
	
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) DecodeRPC(prot *Protocol, senderAddress *net.UDPAddr, connection *net.UDPConn) *Protocol{
	if(prot.Rpc == "PING"){
		return network.PingHandler(prot, senderAddress, connection)
	}else if(prot.Rpc == "NODE_LOOKUP"){
		return network.LookupHandler(prot, senderAddress, connection)
	}else if(prot.Rpc == "NODE_JOIN"){
		return network.JoinHandler(prot, senderAddress, connection)
	}else if(prot.Rpc == "NODE_VALUE"){
		return nil
	}else if(prot.Rpc == "NODE_STORE"){
		return nil
	}else{
		fmt.Println("ERROR. RPC TYPE COULD NOT BE FOUND")
		return nil
	}
}

func (network *Network) JoinHandler(prot *Protocol, responseAddr *net.UDPAddr, connection *net.UDPConn) *Protocol{
	fmt.Println("Inside join handler")
	if(prot.Message == "NODE_JOIN_SENT"){

		//Adds contact
		sendContact := Contact{}
		json.Unmarshal(prot.Data[:len(prot.Data)], &sendContact)
		network.routingTable.AddContact(sendContact) 

		fmt.Println(sendContact.ID)

		
		for i := 0; i < 160 ; i++ {
			if(network.routingTable.buckets[i].Len() > 0){
				fmt.Println(network.routingTable.buckets[i].list.Front())
			}
			
		}
		
		//Respond with my own contact
		meContact := ContactToByteArray(&network.routingTable.me)
		joinProtocolResponse := CreateProtocol("NODE_JOIN",nil,"",meContact,"NODE_JOIN_RESPONSE")
		_, e := connection.WriteToUDP(joinProtocolResponse, responseAddr)
		if e != nil{
			fmt.Println("JoinHandler ERROR")
			fmt.Println(e)
		}
		return prot
	}else if(prot.Message == "NODE_JOIN_RESPONSE"){
		return prot
	}
	return nil
}

func (network *Network) LookupHandler(prot *Protocol, responseAddr *net.UDPAddr, connection *net.UDPConn) *Protocol{
	if(prot.Message == "NODE_LOOKUP_SENT"){
		targetContact := Contact{}
		json.Unmarshal(prot.Data[:len(prot.Data)], &targetContact)
		closestContactsArray := network.routingTable.FindClosestContacts(targetContact.ID, 3)
		fmt.Println(len(closestContactsArray))
		lookupProtocolResponse := CreateProtocol("NODE_LOOKUP",closestContactsArray,"",prot.Data,"NODE_LOOKUP_RESPONSE")
		_, e := connection.WriteToUDP(lookupProtocolResponse, responseAddr)
		if e != nil{
			fmt.Println("LookupHandler ERROR")
			fmt.Println(e)
		}
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
		fmt.Println(prot.Message)
		pingResponseRPC := CreateProtocol("PING",nil,"",nil,"PING_RESPONSE")
		_, e := connection.WriteToUDP(pingResponseRPC, responseAddr)
		if e != nil {
			fmt.Println("PingHandler ERROR")
			fmt.Println(e)
		}
		return prot
	}else if(prot.Message == "PING_RESPONSE"){
		fmt.Println(prot.Message)
		return prot
	}
	return nil
}

func ContactToByteArray(contact *Contact) []byte {
	contactByteArray, err := json.Marshal(contact)
	if err != nil{
		fmt.Println(err)
		fmt.Println("FAIL IN CONVERSION")
		return nil
	}
	return contactByteArray
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

/*
	Change the InitNetwork to the newly defined struct. Remove contact and add routingtable
*/

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
