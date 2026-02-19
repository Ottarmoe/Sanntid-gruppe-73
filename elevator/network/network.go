package network

import (
	"elevator/networkLow"
	"fmt"
	"time"
)

func TestNodeCommunication(id string) error {
	go periodicSendingOfId(id)

	for {
		data, err := networkLow.Receive()
		if err != nil {
			fmt.Println("Receive error:", err)
			continue // keep listening
		}

		networkLow.PrintMessage(data)
	}
}

func periodicSendingOfId(id string) {
	for i := 0; ; i++ {
		data := []byte(fmt.Sprintf("Hello %d from %s", i, id))

		if err := networkLow.Send(data); err != nil {
			fmt.Println("Send error:", err)
		}

		time.Sleep(1 * time.Second)
	}
}