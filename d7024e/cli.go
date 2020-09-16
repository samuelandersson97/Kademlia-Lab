package d7024e

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

func scanInput() {
	reader := bufio.NewReader(os.Stdin)
	input := strings.Split(reader.ReadString('\n'), " ")
	handleInput(input)
}

func handleInput(s []string) {
	operation := s[0]
	if operation == "ping"{
		testContact := NewContact(util.NewRandomKademliaID(), "127.0.0.1")
		testNetwork := CreateNetwork(&testContact)
	}
	
	if operation == "node"{
		if s[1] == "lookup"{

		}
		if s[1] == "join"{

		}
		else{
			fmt.Println("Incorrect command!")
		}
	}
	
	if operation == "put"{
		
	}
	
	if operation == "get"{
		
	}
	
	if operation == "exit"{
		os.Exit()
	}
	
	if operation == "help"{
		cmdString := "\n
		Available commands:\n
		\n
		"
		fmt.Println(cmdString)
	}
	
	else{
		fmt.Println("Incorrect command!")
	}
}

func cliGreeting() {
	greeting := "\n
	##################################################
	#     WELCOME TO THE KADEMLIA CLI.               #
	# PLEASE ENTER YOUR COMMANDS IN THE TERMINAL     #
	# WRITE 'help' FOR A LIST OF AVAILABLE COMMANDS  #
	##################################################
	\n
	"
	fmt.Println(greeting)
}

