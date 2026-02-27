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
	addr.Port = 20001
	addr.IP = net.IPv4(0, 0, 0, 0)

	
	var toAddr net.UDPAddr
	toAddr.Port = 20000
	//toAddr.IP = net.IPv4(10, 0, 10, 111) //use this when running server_wfh.exe
	toAddr.IP = net.IPv4(172, 17, 146, 72) //use this when running server_wfh

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
}
