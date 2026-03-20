package networkLow

import (
	. "elevator/elevatorConstants"
	"net"
	// "syscall"
)

var conn *net.UDPConn

var broadcastAddr = &net.UDPAddr{
	IP:   net.IPv4bcast,
	Port: CommunicationPort,
}

var broadcastReceiveAddr = &net.UDPAddr{
	IP:   net.IPv4zero,
	Port: CommunicationPort,
}

func Init() error {
	var err error

	conn, err = net.ListenUDP("udp", broadcastReceiveAddr)
	if err != nil {
		return err
	}

	// rawConn, err := conn.SyscallConn()
	// if err != nil {
	// 	return err
	// }

	// err = rawConn.Control(func(fd uintptr) {
	// 	syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

func Send(data []byte) error {
	_, err := conn.WriteToUDP(data, broadcastAddr)
	return err
}

func Receive() ([]byte, error) {
	buf := make([]byte, 1024)

	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}
