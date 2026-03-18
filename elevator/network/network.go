package network

import (
	"elevator/networkLow"
	. "elevator/stateTypes"
	"fmt"
	"time"

	. "elevator/elevatorConstants"
	"encoding/json"
	"log"
)

func NetworkSender(netMessageToNetworkSender <-chan NetMessage, pokeState chan<- struct{}) {
	timeToSend := time.NewTicker(time.Duration(BroadcastRate * float64(time.Second)))

	netMessage := <-netMessageToNetworkSender

	for {
		select {
		case netMessage = <-netMessageToNetworkSender:

		case <-timeToSend.C:
			pokeState <- struct{}{}
			data, err := json.Marshal(netMessage)
			//fmt.Println(string(data))
			if err != nil {
				log.Println("send jsonmarshal error:", err)
			}
			err = networkLow.Send(data)
			if err != nil {
				log.Println("send error:", err)
			}
		}
	}
}

func NetworkReceiver(netMessageToState chan<- NetMessage, netErrorToState chan<- NetErrorNotification) {
	var prevNetMessages [NumElevators]NetMessage

	receiveMessage := make(chan NetMessage)
	go receiver(receiveMessage)

	var NetError [NumElevators]bool
	var ReceivedDuringInterval [NumElevators]int
	var ReceivedIntervals [NumElevators]int
	for i := 0; i < NumElevators; i++ {
		NetError[i] = true
		ReceivedDuringInterval[i] = 0
		ReceivedIntervals[i] = 0
	}
	timeout := make(chan int)
	var resetTimer [NumElevators]chan struct{}
	for i := 0; i < NumElevators; i++ {
		if i == ID() {
			continue
		}
		resetTimer[i] = make(chan struct{}, 1)
		go timoutNotifier(i, timeout, resetTimer[i])
	}

	for {
		select {
		case netMessage := <-receiveMessage:
			//Handle neterror
			if NetError[netMessage.ID] {
				ReceivedDuringInterval[netMessage.ID]++
				if ReceivedDuringInterval[netMessage.ID] < MessagesNeededWithinInterval {
					continue
				} else {
					ReceivedIntervals[netMessage.ID]++
					ReceivedDuringInterval[netMessage.ID] = 0
					if ReceivedIntervals[netMessage.ID] < IntervalsNeeded {
						continue
					} else {
						NetError[netMessage.ID] = false
						netErrorToState <- NetErrorNotification{ID: netMessage.ID, NetError: false}
					}
				}
			}
			resetTimer[netMessage.ID] <- struct{}{} //concider making non-blocking...seems to work for now

			//Avoid botherrueing state with duplicate messages
			if netMessage == prevNetMessages[netMessage.ID] {
				//fmt.Print("c")
				continue
			}
			fmt.Print("unique message")
			fmt.Println(netMessage)

			netMessageToState <- netMessage
			prevNetMessages[netMessage.ID] = netMessage

		case id := <-timeout:
			if !NetError[id] {
				NetError[id] = true
				netErrorToState <- NetErrorNotification{ID: id, NetError: true}
				ReceivedDuringInterval[id] = 0
				ReceivedIntervals[id] = 0
				resetTimer[id] <- struct{}{}
			} else {
				ReceivedDuringInterval[id] = 0
				ReceivedIntervals[id] = 0
				resetTimer[id] <- struct{}{}
			}
		}
	}
}

func timoutNotifier(id int, timeout chan int, resetTimer chan struct{}) {
	timer := time.NewTimer(time.Duration(NetErrorTimerLength * float64(time.Second)))
	for {
		select {
		case <-timer.C:
			timeout <- id
		case <-resetTimer:
			if !timer.Stop() {
				// Timer already fired, drain the channel to avoid spurious tick
				//Check documention if this is still needed! https://pkg.go.dev/time#NewTimer
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(time.Duration(NetErrorTimerLength * float64(time.Second)))
		}
	}
}

func receiver(receiveMessage chan<- NetMessage) {
	// last := time.Now() //For debugging

	for {
		data, err := networkLow.Receive()
		if err != nil {
			// log.Println("receive  error:", err)
			continue
		}
		var netMessage NetMessage
		err = json.Unmarshal(data, &netMessage)
		//fmt.Println("unmarshaled", netMessage)
		if err != nil {
			log.Println("receive json unmarshal error:", err)
			continue
		}
		if netMessage.ID == ID() {
			continue
		}
		//fmt.Println("recieved", netMessage)

		// For debugging
		// now := time.Now()
		// fmt.Println(now.Sub(last))
		// last = now

		receiveMessage <- netMessage
	}

}
