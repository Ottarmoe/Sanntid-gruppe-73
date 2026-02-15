package network

func Transmitter(){
	conn := conn.DialBroadcastUDP(port)

	conn.WriteTo(ttj, addr)
}