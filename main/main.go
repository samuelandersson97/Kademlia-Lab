package main

import (
	"d7024e"
)

func main() {
	d7024e.CliGreeting()
	for{
		//d7024e.ScanInput()
		var ping []string
		ping = append(ping, "ping","boi")
		d7024e.HandleInput(ping)
	}
}
