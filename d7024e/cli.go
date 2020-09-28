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

func ScanInput(me Contact) {
	reader := bufio.NewReader(os.Stdin)
	readValue, err := reader.ReadString('\n')
	inputString := strings.Split(readValue, "\n")
	if err != nil {
		fmt.Println(err)
	}
	input := strings.Split(inputString[0], " ")
	HandleInput(input, me)
}

/*
	Should add support for 'node lookup', 'node join', 'put' and 'get' when we are finished creating the support for these operations.

	Note that the 'ping' operation should be changed since we are creating networks and contacts there as it stands just to test the operation.

*/

func HandleInput(s []string, me Contact) {
	operation := s[0]
	if operation == "ping"{
		testContact := NewContact(NewRandomKademliaID(), s[1])
		testNetwork := InitNetwork(&testContact)
		testNetwork.SendPingMessage(&testContact)
	}else if operation == "node"{
		if s[1] == "lookup"{
			testKad := InitKademlia(me)
			testContact := NewContact(NewRandomKademliaID(), s[1])
			testNetwork := InitNetwork(&testContact)
			testKad.LookupContact(&testContact)
		}else if s[1] == "join"{

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

