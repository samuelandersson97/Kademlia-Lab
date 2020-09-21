package main

import (
	"d7024e"
)

func main() {
	myIP := d7024e.GetOutboundIP()
	go d7024e.Listen(myIP, 1111)
	d7024e.CliGreeting()
	for{
		d7024e.ScanInput()
	}
}
