package main

import (
	"d7024e"
	"net"
	"fmt"
)

func GetOutboundIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        fmt.Println(err)
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP.String()
}

func main() {
	myIP := GetOutboundIP()
	go d7024e.Listen(myIP, 1111)
	d7024e.CliGreeting()
	for{
		d7024e.ScanInput()
	}
}
