package main

import (
	"d7024e"
	"fmt"
	"time"
)
/*
	Create a new instance of the kademlia struct.
	Needs a routing table and the 'me contact'
*/

func main() {
	myIP := d7024e.GetOutboundIP()
	randId := d7024e.NewRandomKademliaID(time.Now().UnixNano())
	fmt.Println(randId.String())
	me := d7024e.NewContact(randId,myIP)
	rt := d7024e.NewRoutingTable(me)
	network := d7024e.InitNetwork(rt)
	kad := d7024e.InitKademlia(network)
	go network.Listen(myIP, 1111)
	d7024e.CliGreeting()
	for{
		d7024e.ScanInput(kad)
	}
} 
