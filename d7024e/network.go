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
	//Testing for ping only, should be more than one contact
	contact *Contact
	// Routing table, not contact

}

type Protocol struct {
	Rpc string
	Contacts []*Contact
	Hash string
	Data []byte
	Message string
}


/*
	Extract the sent message and create different fucntions that handles different types of messages (PING, FIND_NODE, etc...)
*/

func Listen(ip string, port int) {
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
		DecodeRPC(&responseProtocol, senderAddress, c)
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
		DecodeRPC(&receivedPing, senderAddress, c)
	}
	
}

func CreateProtocol(rpcToSend string, contactsArr []*Contact, hashToSend string, dataToSend []byte, messageToSend string) []byte{
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

func (network *Network) SendFindContactMessage(contact *Contact) {
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
	lookupMessage := CreateProtocol("NODE_LOOKUP", nil, "", nil, "NODE_LOOKUP_SENT")
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
	}else{
		json.Unmarshal(responseBuffer[:size], &receivedLookup)
		DecodeRPC(&receivedLookup, senderAddress, c)
	}
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func DecodeRPC(prot *Protocol, senderAddress *net.UDPAddr, connection *net.UDPConn){
	if(prot.Rpc == "PING"){
		PingHandler(prot, senderAddress, connection)
	}else if(prot.Rpc == "NODE_LOOKUP"){
		fmt.Println("Inside decoder")
		LookupHandler(prot)
	}else if(prot.Rpc == "NODE_VALUE"){

	}else if(prot.Rpc == "STORE"){

	}else{
		fmt.Println("ERROR. RPC TYPE COULD NOT BE FOUND")
	}
}

func LookupHandler(prot *Protocol){
	fmt.Println("Inside handler")
}

func PingHandler(prot *Protocol, responseAddr *net.UDPAddr, connection *net.UDPConn){
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
	}else if(prot.Message == "PING_RESPONSE"){
		fmt.Println(prot.Message)
	}
}

/*
	Change the InitNetwork to the newly defined struct. Remove contact and add routingtable
*/

func InitNetwork(contact *Contact) *Network{
	network := &Network{}
	network.contact = contact
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
