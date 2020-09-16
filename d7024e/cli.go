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

	}
	
	else if operation == "node"{
		if s[1] == "lookup"{

		}
		else if s[1] == "join"{

		}
		else{
			fmt.Println("Incorrect command!")
		}
	}
	
	else if operation == "put"{
		
	}
	
	else if operation == "get"{
		
	}
	
	else if operation == "exit"{
		os.Exit()
	}
	
	else if operation == "help"{
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
	#		   WELCOME TO THE KADEMLIA CLI.			 #
	#	PLEASE ENTER YOUR COMMANDS IN THE TERMINAL   #
	# WRITE 'help' FOR A LIST OF AVAILABLE COMMANDS  #
	##################################################
	\n
	"
	fmt.Println(greeting)
}

