package main

import (
	"fmt"
	"net"
	"time"
)

func receiver(conn *net.UDPConn){
	buf := make([]byte, 2048)

		for {
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("Recive error:", err)
				continue
			}
			fmt.Printf("from %s: %s\n", addr.String(), string(buf[:n]))
		}

}

func sender(conn *net.UDPConn, serverAddr *net.UDPAddr){
	for i := 0; ; i++ {
		msg := fmt.Sprintf("Hello %d", i)
		_, err := conn.WriteToUDP([]byte(msg), serverAddr)
		if err != nil {
			fmt.Println("Send error:", err)
		}
		time.Sleep(1000000000)
	}
}


func main() {	
	
	var addr net.UDPAddr
	addr.Port = 20002
	addr.IP = net.IPv4(0, 0, 0, 0)

	
	var toAddr net.UDPAddr
	toAddr.Port = 20002
	toAddr.IP = net.IPv4(10, 100, 23, 11)

	recvConn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer recvConn.Close()


	sendConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(err)
	}
	defer sendConn.Close()

	go receiver(recvConn)
	go sender(sendConn, &toAddr)

	select {}

	// sock, err := net.DialUDP("udp", &sendAddr, &toAddr)
	// if err != nil {
	// 	panic(err)
	// }
	// var buffer [1024]byte
	// for {
	// 	time.Sleep(200_000_000)

		// if sending directly to a single remote machine:

		// either: set up the socket to use a single remote address
		//sock.connect(addr)
		//sock.send(message)
		// or: set up the remote address when sending
		// sock.Write([]byte("heihei"))

		// n, err := sock.Read(buffer[:])
		// _ = err
		// _ = n

		// fmt.Println(string(buffer[:]))

		// if sending on broadcast:
		// you have to set up the BROADCAST socket option before calling connect / sendTo
		//broadcastIP = #.#.#.255 // First three bytes are from the local IP, or just use 255.255.255.255
		//addr = new InternetAddress(broadcastIP, port)
		//sendSock = new Socket(udp) // UDP, aka SOCK_DGRAM
		//sendSock.setOption(broadcast, true)
		//sendSock.sendTo(message, addr)
}
