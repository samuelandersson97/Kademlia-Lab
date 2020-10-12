package d7024e

import (
	"fmt"
	"testing"
	"time"
)

func TestLookup(t *testing.T) {
	contact0IP := "localhost:4000"
	contact0 := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"),contact0IP)
	rt0 := NewRoutingTable(contact0)
	network0 := InitNetwork(rt0)
	network0.routingTable.AddContact(contact0)
	go network0.Listen(contact0IP)
	<- time.After(3*time.Second)

	contact1IP := "localhost:4001"
	contact1 := NewContact(NewKademliaID("1111111300000000000000000000000000000000"),contact1IP)
	rt1 := NewRoutingTable(contact1)
	network1 := InitNetwork(rt1)
	network1.routingTable.AddContact(contact1)
	go network1.Listen(contact1IP)
	<- time.After(3*time.Second)
	network0.routingTable.AddContact(contact1)

	contact2IP := "localhost:4002"
	contact2 := NewContact(NewKademliaID("1111111200000000000000000000000000000000"),contact2IP)
	rt2 := NewRoutingTable(contact2)
	network2 := InitNetwork(rt2)
	network2.routingTable.AddContact(contact2)
	go network2.Listen(contact2IP)
	<- time.After(3*time.Second)
	network0.routingTable.AddContact(contact2)

	contact3IP := "localhost:4003"
	contact3 := NewContact(NewKademliaID("1111111400000000000000000000000000000000"),contact3IP)
	rt3 := NewRoutingTable(contact3)
	network3 := InitNetwork(rt3)
	network3.routingTable.AddContact(contact3)
	go network3.Listen(contact3IP)
	<- time.After(3*time.Second)
	network0.routingTable.AddContact(contact3)

	kad := InitKademlia(network0)
	contacts := kad.LookupContact(&contact3)
	for _, c := range contacts{
		fmt.Println(c.Address)
	}
	
}

func TestStore(t *testing.T) {
	contact0IP := "localhost:4005"
	contact0 := NewContact(NewKademliaID("ffffffff80000000000000000000000000000000"),contact0IP)
	rt0 := NewRoutingTable(contact0)
	network0 := InitNetwork(rt0)
	network0.routingTable.AddContact(contact0)
	go network0.Listen(contact0IP)
	<- time.After(3*time.Second)

	contact1IP := "localhost:4006"
	contact1 := NewContact(NewKademliaID("1111111390000000000000000000000000000000"),contact1IP)
	rt1 := NewRoutingTable(contact1)
	network1 := InitNetwork(rt1)
	network1.routingTable.AddContact(contact1)
	go network1.Listen(contact1IP)
	<- time.After(3*time.Second)
	network0.routingTable.AddContact(contact1)

	kad := InitKademlia(network0)

	Data := "TestData"
	kad.Store([]byte(Data))
	
}

func TestGet(t *testing.T) {
	contact0IP := "localhost:4007"
	contact0 := NewContact(NewKademliaID("ffffffff81000000000000000000000000000000"),contact0IP)
	rt0 := NewRoutingTable(contact0)
	network0 := InitNetwork(rt0)
	network0.routingTable.AddContact(contact0)
	go network0.Listen(contact0IP)
	<- time.After(3*time.Second)

	contact1IP := "localhost:4008"
	contact1 := NewContact(NewKademliaID("1111111391000000000000000000000000000000"),contact1IP)
	rt1 := NewRoutingTable(contact1)
	network1 := InitNetwork(rt1)
	network1.routingTable.AddContact(contact1)
	go network1.Listen(contact1IP)
	<- time.After(3*time.Second)
	network0.routingTable.AddContact(contact1)

	kad := InitKademlia(network0)
	
	Data := "TestData"
	kad.Store([]byte(Data))
	hash := kad.network.hashtable[0].Key.String()
	kad.LookupData(hash)
	contact2IP := "localhost:4019"
	contact2 := NewContact(NewKademliaID("1111111392000000000000000000000000000000"),contact2IP)
	rt2 := NewRoutingTable(contact2)
	network2 := InitNetwork(rt2)
	network2.routingTable.AddContact(contact2)
	go network2.Listen(contact2IP)
	<- time.After(3*time.Second)
	network2.routingTable.AddContact(contact0)
	kad2 := InitKademlia(network2)
	kad2.LookupData(hash)
}

func TestJoin(t *testing.T) {
	contact0IP := "localhost:4027"
	contact0 := NewContact(NewKademliaID("ffffffff81100000000000000000000000000000"),contact0IP)
	rt0 := NewRoutingTable(contact0)
	network0 := InitNetwork(rt0)
	network0.routingTable.AddContact(contact0)
	go network0.Listen(contact0IP)
	<- time.After(3*time.Second)

	contact1IP := "localhost:4028"
	contact1 := NewContact(NewKademliaID("1111111391100000000000000000000000000000"),contact1IP)
	rt1 := NewRoutingTable(contact1)
	network1 := InitNetwork(rt1)
	network1.routingTable.AddContact(contact1)
	go network1.Listen(contact1IP)
	<- time.After(3*time.Second)
	network0.routingTable.AddContact(contact1)

	kad := InitKademlia(network0)
	
	kad.NodeJoin(contact1.Address)
}

func TestPing(t *testing.T) {
	contact0IP := "localhost:4029"
	contact0 := NewContact(NewKademliaID("ffffffff81220000000000000000000000000000"),contact0IP)
	rt0 := NewRoutingTable(contact0)
	network0 := InitNetwork(rt0)
	network0.routingTable.AddContact(contact0)
	go network0.Listen(contact0IP)
	<- time.After(3*time.Second)

	contact1IP := "localhost:4030"
	contact1 := NewContact(NewKademliaID("1111111392120000000000000000000000000000"),contact1IP)
	rt1 := NewRoutingTable(contact1)
	network1 := InitNetwork(rt1)
	network1.routingTable.AddContact(contact1)
	go network1.Listen(contact1IP)
	<- time.After(3*time.Second)
	network0.routingTable.AddContact(contact1)

	kad := InitKademlia(network0)
	
	kad.network.SendPingMessage(&contact1)
}