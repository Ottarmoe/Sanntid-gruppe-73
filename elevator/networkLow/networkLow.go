package networkLow

import (
	"fmt"
	"net"
	"syscall"
	"context"
	"golang.org/x/sys/unix"
)

var conn *net.UDPConn

var broadcastAddr = &net.UDPAddr{
	IP:   net.IPv4bcast, // 255.255.255.255
	Port: 30000,
}

func Init() error {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var controlErr error
			err := c.Control(func(fd uintptr) {
				if err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
					controlErr = err
					return
				}
				if err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
					controlErr = err
					return
				}
				if err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_BROADCAST, 1); err != nil {
					controlErr = err
					return
				}
			})
			if err != nil {
				return err
			}
			return controlErr
		},
	}

	pc, err := lc.ListenPacket(context.Background(), "udp", ":30000")
	if err != nil {
		return err
	}

	conn = pc.(*net.UDPConn)
	return nil
}

func Send(data []byte) error {
    _, err := conn.WriteToUDP(data, broadcastAddr)
    return err
}

//returns slice of array with length mathcing the exact length of the message
func Receive() ([]byte, error) {
    buf := make([]byte, 1024)

    n, _, err := conn.ReadFromUDP(buf)
    if err != nil {
        return nil, err
    }

    return buf[:n],nil
}

func PrintMessage(data []byte){
	fmt.Printf("Received %s\n", string(data[:]))

}