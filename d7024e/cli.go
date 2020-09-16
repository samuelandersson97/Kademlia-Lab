package d7024e

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

func ScanInput() {
	reader := bufio.NewReader(os.Stdin)
	readValue, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	input := strings.Split(readValue, " ")
	HandleInput(input)
}

func HandleInput(s []string) {
	var network Network
	operation := s[0]
	if operation == "ping"{
		testContact := NewContact(NewRandomKademliaID(), "10.0.1.22")
		testNetwork := network.InitNetwork(&testContact)
		testNetwork.SendPingMessage(&testContact)
	}
	
	if operation == "node"{
		if s[1] == "lookup"{

		}
		if s[1] == "join"{

		}else{
			fmt.Println("Incorrect command!")
		}
	}
	
	if operation == "put"{
		
	}
	
	if operation == "get"{
		
	}
	
	if operation == "exit"{
		os.Exit(0)
	}
	
	if operation == "help"{
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

