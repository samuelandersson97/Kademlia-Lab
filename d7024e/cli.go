package d7024e

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

/*
	Split the input string.
	Call handle input with the correct string array.
*/

func ScanInput(kad *Kademlia) {
	reader := bufio.NewReader(os.Stdin)
	readValue, err := reader.ReadString('\n')
	inputString := strings.Split(readValue, "\n")
	if err != nil {
		fmt.Println(err)
	}
	input := strings.Split(inputString[0], " ")
	HandleInput(input, kad)
}

/*
	Interpret the command and call the correct function in Kademlia. 
*/

func HandleInput(s []string,  kad *Kademlia) {
	operation := s[0]
	if operation == "ping"{
		/*
			Pinging the ip given. 
		*/
			contact := NewContact(NewRandomKademliaID(), s[1]) // Dummy contact with the correct ip and a random kademlia id
			kad.network.SendPingMessage(&contact)
	}else if operation == "node"{
		if s[1] == "lookup"{
			/*
				Look up the node with the given ID
			*/
			contactNoIp := NewContact(NewKademliaID(s[2]), "")// Dummy contact with the correct id and no ip
			kad.network.routingTable.AddContact(contactNoIp)
			kad.LookupContact(&contactNoIp)
		}else if s[1] == "join"{
			/*
				Joins the node given
			*/
			kad.NodeJoin(s[2])
		
		}else{
			// The command was wrong!
			fmt.Println("Incorrect command!")
		}
	}else if operation == "put"{
		// Store the data on this node and nodes with id close to the key
		kad.Store([]byte(s[1]))
	}else if operation == "get"{
		// Search for the data on this node or nodes with id close to the key
		s := kad.LookupData(s[1])
		fmt.Println("Data found: "+s)
	}else if operation == "show"{
		// Displays the data stored
		for _, d := range kad.network.hashtable {
			fmt.Println(string(d.Data))
		}
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
	/*
		This is only a greeting
	*/
	greeting := "WELCOME TO THE KADEMLIA CLI.\n"
	fmt.Println(greeting)
}

