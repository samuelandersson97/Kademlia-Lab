package main

import (
	"d7024e"
)
/*
	Create a new instance of the kademlia struct.
	Needs a routing table and the 'me contact'


	Check if this node is start node, give the node an id in that case.
	Otherwise Id is assigned during join sequence. 
*/

func main() {
	myIP := d7024e.GetOutboundIP()
	me := d7024e.NewContact(d7024e.NewRandomKademliaID(),myIP)
	rt := d7024e.NewRoutingTable(me)
	network := d7024e.InitNetwork(rt)
	kad := d7024e.InitKademlia(network)
	go network.Listen(myIP, 1111)
	d7024e.CliGreeting()
	for{
		d7024e.ScanInput(kad)
	}
} 
