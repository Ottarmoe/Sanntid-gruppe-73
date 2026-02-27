package network

import (
	"elevator/networkLow"
	"fmt"
	"time"
)

func TestNodeCommunication(id int) error {
	go periodicSendingOfId(id)

	for {
		data, err := networkLow.Receive()
		if err != nil {
			fmt.Println("Receive error:", err)
			continue
		}

		networkLow.PrintMessage(data)
	}
}

func periodicSendingOfId(id int) {
	for i := 0; ; i++ {
		data := []byte(fmt.Sprintf("Hello %d from %d", i, id))

		err := networkLow.Send(data)
		if err != nil {
			fmt.Println("Send error:", err)
		}

		time.Sleep(1 * time.Second)
	}
}