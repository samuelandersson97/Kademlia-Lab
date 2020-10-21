# Kademlia-Lab
Lab in D7024E

To start up the nodes, run:

> ./startScript.sh

Then you can attach to any of the nodes and run the commands below


# Commands:

> put "data"
  
stores the data internally and to the closest nodes

> get "data"
  
gets the data from either itself of from any other node that has the data

> node join "address"
  
joins the node with that address

> ping "address"

pings the node with that address

> show

prints the content of the nodes hashtable
