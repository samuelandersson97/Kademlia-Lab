package d7024e

import(
	"net"
	"strconv"
	"fmt"
)

type Network struct {
	//Testing for ping only, should be more than one contact
	contact *Contact
}

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
	for{
		//Adds the message from the UDP-channel in the message-buffer. Returns the size of the message and the adress of the sender
		size, senderAddress, err := c.ReadFromUDP(messageBuffer)
		if err != nil {
			fmt.Println("LISTEN ERROR: 3")
			fmt.Println(err)
		}
		fmt.Println("This was sent from "+ senderAddress.String() +": "+string(messageBuffer[0:size-1])+"\n")
		response := []byte("Det är klart att jag kan!")
		_, e := c.WriteToUDP(response, senderAddress)
		if e != nil {
			fmt.Println("LISTEN ERROR: 4")
			fmt.Println(e)
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	//Returns an address of the UDP end point. 'udp4' indicates that only IPv4-addresses are being resolved
	fmt.Println("This is contact ip: " + contact.Address)
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
	defer c.Close()
	message := []byte("Halloj! Kan du snälla besvara mig?")
	_, e := c.Write(message)
	if e != nil {
		fmt.Println("SEND ERROR: 3")
		fmt.Println(err)
	}
	messageBuffer := make([]byte, 8192)
	size, senderAddress, err := c.ReadFromUDP(messageBuffer)
	if err != nil {
		fmt.Println("SEND ERROR: 4")
		fmt.Println(err)
	}
	fmt.Println("RESPONSE: "+ senderAddress.String() +": "+string(messageBuffer[0:size-1])+"\n")
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