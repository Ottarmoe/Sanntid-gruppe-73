package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	// the address we are listening for messages on
	// we have no choice in IP, so use 0.0.0.0, INADDR_ANY, or leave the IP field empty
	var addr net.UDPAddr

	// the port should be whatever the sender sends to
	addr.Port = 30000
	addr.IP = net.IPv4(0, 0, 0, 0)

	// a socket that plugs our program to the network. This is the "portal" to the outside world
	// alternate names: conn
	// UDP is sometimes called SOCK_DGRAM. You will sometimes also find UDPSocket or UDPConn as separate types
	//recvSock = new Socket(udp)
	conn, err := net.ListenUDP("udp", &addr)
	_ = err

	// bind the address we want to use to the socket
	//recvSock.bind(addr)

	// a buffer where the received network data is stored
	var buffer [1024]byte

	// an empty address that will be filled with info about who sent the data
	//var fromWho net.UDPAddr

	for {
		// clear buffer (or just create a new one)
		for index := range buffer {
			buffer[index] = 0
		}
		// receive data on the socket
		// fromWho will be modified by ref here. Or it's a return value. Depends.
		// receive-like functions return the number of bytes received
		// alternate names: read, readFrom
		n, fromWho, err := conn.ReadFromUDP(buffer[:])

		// the buffer just contains a bunch of bytes, so you may have to explicitly convert it to a string
		//buffer[n] = 0
		fmt.Println(string(buffer[:]))
		fmt.Println(fromWho.Port)

		// optional: filter out messages from ourselves
		//if(fromWho.IP != localIP){
		// do stuff with buffer
		//}
		_ = err
		_ = n
		time.Sleep(50000000)
	}
}
