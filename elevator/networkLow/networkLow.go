package networkLow

import (
	"fmt"
	"net"
	"syscall"
	"context"
	"golang.org/x/sys/unix"
	"sync"
	"math/rand"
	"errors"
	"net/http"
	"strconv"
	. "elevator/elevatorConstants"
)

var conn *net.UDPConn

var broadcastAddr = &net.UDPAddr{
	IP:   net.IPv4bcast, // 255.255.255.255
	Port: 30000,
}


//Simulate packet loss 
//curl "http://localhost:8080/set_loss?prob=0"
var (
	packetLossProb float64 = 0 // 0% initial packet loss
	probMutex      sync.RWMutex   // protects access to packetLossProb
)

func getPacketLossProb() float64 {
	probMutex.RLock()
	defer probMutex.RUnlock()
	return packetLossProb
}

func SetPacketLoss(prob float64) {
	probMutex.Lock()
	defer probMutex.Unlock()
	if prob < 0 {
		prob = 0
	} else if prob > 1 {
		prob = 1
	}
	packetLossProb = prob
	fmt.Printf("Packet loss probability set to %.2f%%\n", packetLossProb*100)
}

var ErrPacketDropped = errors.New("simulated packet loss")

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

	//Simulate packet loss
	port := 8080 + ID()

	http.HandleFunc("/set_loss", func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query().Get("prob")
		if q == "" {
			http.Error(w, "missing prob", 400)
			return
		}

		prob, err := strconv.ParseFloat(q, 64)
		if err != nil {
			http.Error(w, "invalid prob", 400)
			return
		}

		SetPacketLoss(prob)

		fmt.Fprintf(w,
			"Program ID %d packet loss set to %.2f%%\n",
			ID(),
			prob*100,
		)
	})

	addr := fmt.Sprintf(":%d", port)

	fmt.Println("Packet loss simulator running. Initial packet loss 0. Use command: curl http://localhost:8080/set_loss?prob=0.5 in another terminal to change packet loss during runtime. For id != 0, replace 8080 with 8080+id.")

	go http.ListenAndServe(addr, nil)

	return nil
}

func Send(data []byte) error {
	//Simulate packet loss
	if rand.Float64() < getPacketLossProb() {
		// fmt.Println("Packet not sent")
		return nil // drop packet
	}

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

	//Simulate packet loss
	if rand.Float64() < getPacketLossProb() {
		return nil, ErrPacketDropped // drop packet
	}

    return buf[:n],nil
}

func PrintMessage(data []byte){
	fmt.Printf("Received %s\n", string(data[:]))

}