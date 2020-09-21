package main

import (
	"d7024e"
)

func main() {
	d7024e.Listen("127.0.0.1", 1111)
	d7024e.CliGreeting()
	for{
		d7024e.ScanInput()
	}
}
