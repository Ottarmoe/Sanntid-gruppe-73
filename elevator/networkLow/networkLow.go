package networkLow

conn

func init(port int){
	var addr net.UDPAddr
	addr.Port = port
	addr.IP = net.IPv4(0, 0, 0, 0)

	conn = conn.DialBroadcastUDP(port)
	
}




conn, err := net.ListenUDP("udp", &addr)
if err != nil {
	panic(err)
}
defer conn.Close()

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