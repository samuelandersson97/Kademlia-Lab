package d7024e

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

/*
	Should be complete
*/

func ScanInput(kad *Kademlia) {
	reader := bufio.NewReader(os.Stdin)
	readValue, err := reader.ReadString('\n')
	inputString := strings.Split(readValue, "\n")
	if err != nil {
		fmt.Println(err)
	}
	input := strings.Split(inputString[0], " ")
	HandleInput(input, network, kad)
}

/*
	Should add support for 'node lookup', 'node join', 'put' and 'get' when we are finished creating the support for these operations.

	Note that the 'ping' operation should be changed since we are creating networks and contacts there as it stands just to test the operation.

*/

func HandleInput(s []string, network *Network, kad *Kademlia) {
	operation := s[0]
	if operation == "ping"{
		contact := NewContact(NewRandomKademliaID(), s[1])
		kad.network.SendPingMessage(&contact)
	}else if operation == "node"{
		if s[1] == "lookup"{
			contact := NewContact(NewRandomKademliaID(), s[2])
			kad.network.routingTable.AddContact(contact)
			kad.LookupContact(&contact)
		}else if s[1] == "join"{
			kad.NodeJoin(s[2])
		/*
			1.	Ip address supplied to the node we are joining.
			2.	Random id is supplied to this node.   
			3.	K-bucket is initialised with the node that we first know, collect information about this node. 
			4.	Lookup on itself to gain close nodes and the routing table is then updated in this function. 
			5.	Done!
		*/
		}else{
			fmt.Println("Incorrect command!")
		}
	}else if operation == "put"{
		
	}else if operation == "get"{
		
	}else if operation == "exit"{
		os.Exit(0)
	}else if operation == "help"{
		cmdString := "\n Available commands:\n"
		fmt.Println(cmdString)
	}else{
		fmt.Println("Incorrect command!")
	}
}

func CliGreeting() {
	greeting := "WELCOME TO THE KADEMLIA CLI.\n"
	fmt.Println(greeting)
}

