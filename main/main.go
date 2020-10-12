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
	d7024e.CliGreeting()
	myIP := d7024e.GetOutboundIP()+"1111"
	
	me := d7024e.NewContact(d7024e.NewRandomKademliaID(),myIP)
	rt := d7024e.NewRoutingTable(me)
	network := d7024e.InitNetwork(rt)
	kad := d7024e.InitKademlia(network)
	go network.Listen(myIP)
	/*if(!(myIP == "172.1.17.2")){
		kad.NodeJoin("172.1.17.2")
	}*/
	for{
		d7024e.ScanInput(kad)
	}
} 
