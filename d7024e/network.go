package d7024e

import(
	"net"
	"strconv"
	"time"
)

type Network struct {
	//Testing for ping only, should be more than one contact
	contact *Contact
}

//type Network struct {
//	rTable *RoutingTable
//}

func Listen(ip string, port int) {
	adrPort := ip+":"+port
	//Returns an address of the UDP end point. 'udp4' indicates that only IPv4-addresses are being resolved
	udpEndPoint, err := net.ResolveUDPAddr("udp4",adrPort)
	if err != nil {
		fmt.Println(err)
	}
	//Listens for packets on the (ONLY!!) LOCAL network. 'udp4' indicates that only IPv4-addresses are taken into account when it comes to listening for packets, returns a connection
	c, err := net.ListenUDP("udp4", udpEndPoint)
	if err != nil {
		fmt.Println(err)
	}
	defer c.close()
	//creates buffer with maximum length of 512
	messageBuffer := make([]byte, 512)
	for{
		//Adds the message from the UDP-channel in the message-buffer. Returns the size of the message and the adress of the sender
		size, senderAddress, err = c.ReadFromUDP(messageBuffer)
		fmt.Println("This was sent from "+ senderAddress.String() +": "+string(messageBuffer[0:size-1])+"\n")
		response = []byte("This is "+adrPort+"'s response: JAARRÅ\n")
		_, err = c.WriteToUPD(response, senderAddress)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	go Listen(contact.Address, 1000) <- time.After(2*time.Second) //Kicks of a new thread and executes the Listen function on it for two seconds
	//Returns an address of the UDP end point. 'udp4' indicates that only IPv4-addresses are being resolved
	udpEndPoint, err := net.ResolveUDPAddr("udp4",contact.Address)
	if err != nil {
		fmt.Println(err)
	}
	// Starts up a UDP-connection to the resolved UDP-address 
	c, err := net.DialUDP("udp4", udpEndPoint)
	if err != nil {
		fmt.Println(err)
	}
	defer c.close()
	message = byte[]("Halloj")
	_, err = c.Write(message)
	if err != nil {
		fmt.Println(err)
	}
	messageBuffer := make([]byte, 512)
	size, senderAddress, err := c.ReadFromUDP(messageBuffer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("This was sent from "+ senderAddress.String() +": "+string(messageBuffer[0:size-1])+"\n"()
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) CreateNetwork(contact *Contact) *Network {
	network := Network{}
	network.contact = contact
	return network
}
