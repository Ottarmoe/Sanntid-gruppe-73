package networkLow

import (
	"fmt"
	"net"
	"syscall"
)

var conn *net.UDPConn

func Init() error {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var controlErr error
			err := c.Control(func(fd uintptr) {
				// Allow multiple processes to bind same port
				if err := syscall.SetsockoptInt(int(fd),
					syscall.SOL_SOCKET,
					syscall.SO_REUSEADDR, 1); err != nil {
					controlErr = err
					return
				}

				if err := syscall.SetsockoptInt(int(fd),
					syscall.SOL_SOCKET,
					syscall.SO_REUSEPORT, 1); err != nil {
					controlErr = err
					return
				}

				// Allow broadcast sending
				if err := syscall.SetsockoptInt(int(fd),
					syscall.SOL_SOCKET,
					syscall.SO_BROADCAST, 1); err != nil {
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

func Send(msg []byte, to *net.UDPAddr) error {
    _, err := conn.WriteToUDP(msg, to)
    return err
}

func Receive(buf []byte) (int, *net.UDPAddr, error) {
    return conn.ReadFromUDP(buf)
}

func PrintMessage(buf []byte, n int, addr *net.UDPAddr){
	fmt.Printf("from %s: %s\n", addr.String(), string(buf[:n]))
}



// func receiver(conn *net.UDPConn){
// 	buf := make([]byte, 2048)

// 		for {
// 			n, addr, err := conn.ReadFromUDP(buf)
// 			if err != nil {
// 				fmt.Println("Recive error:", err)
// 				continue
// 			}
// 			fmt.Printf("from %s: %s\n", addr.String(), string(buf[:n]))
// 		}

// }

// func sender(conn *net.UDPConn, serverAddr *net.UDPAddr){
// 	for i := 0; ; i++ {
// 		msg := fmt.Sprintf("Hello %d", i)
// 		_, err := conn.WriteToUDP([]byte(msg), serverAddr)
// 		if err != nil {
// 			fmt.Println("Send error:", err)
// 		}
// 		time.Sleep(1000000000)
// 	}
// }