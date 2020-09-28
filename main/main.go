package main

import (
	"d7024e"
)
/*
	Create a new instance of the kademlia struct.
	Needs a routing table and the 'me contact'
*/

func main() {
	myIP := d7024e.GetOutboundIP()
	me := d7024e.NewContact(NewRandomKademliaID(),myIP,nil)
	go d7024e.Listen(myIP, 1111)
	d7024e.CliGreeting()
	for{
		d7024e.ScanInput(me)
	}
}
