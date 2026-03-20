package network

import (
	. "elevator/elevatorConstants"
	"elevator/networkLow"
	. "elevator/stateTypes"
	"encoding/json"
	"log"
	"time"
)

// Stores a network message, which contains the newest state information.
// The routine periodically broadcasts the stored message on the network.
func NetworkSender(netMessageToNetworkSender <-chan NetMessage) {
	timeToSend := time.NewTicker(time.Duration(BroadcastRate * float64(time.Second)))

	//No periodic sending before the first netmessage has been comunnicated by state
	netMessage := <-netMessageToNetworkSender

	for {
		select {
		case netMessage = <-netMessageToNetworkSender:
			//Store newest netMessage

		case <-timeToSend.C:
			data, err := json.Marshal(netMessage)
			if err != nil {
				log.Println("Send json marshal error:", err)
				continue
			}
			err = networkLow.Send(data)
			if err != nil {
				log.Println("Send error:", err)
			}
		}
	}
}

// Receives network messages, filters out duplicate messages and tracks how often
// a message is received from each elevator. If no message is received within a given interval,
// the elevator is marked as neterror. To remove neterror, a certain number of messages must be received
// during an interval, for x intervals in a row. Unique netmessages and neterror-change is sent to state.
func NetworkReceiver(netMessageToState chan<- NetMessage, netErrorToState chan<- NetErrorNotification) {
	var prevNetMessages [NumElevators]NetMessage

	receiveMessage := make(chan NetMessage)
	go receiver(receiveMessage)

	//Track network error state and received messages per interval, nr. of intervals
	var NetError [NumElevators]bool
	var ReceivedDuringInterval [NumElevators]int
	var ReceivedIntervals [NumElevators]int
	for i := 0; i < NumElevators; i++ {
		NetError[i] = true
		ReceivedDuringInterval[i] = 0
		ReceivedIntervals[i] = 0
	}

	// Start a timeout notifier for each other elevator
	timeout := make(chan int)
	var resetTimer [NumElevators]chan struct{}
	for i := 0; i < NumElevators; i++ {
		if i == ID() {
			continue
		}
		resetTimer[i] = make(chan struct{}, 1)
		go timeoutNotifier(i, timeout, resetTimer[i])
	}

	for {
		select {
		case netMessage := <-receiveMessage:
			//If in neterror: handle getting back online
			if NetError[netMessage.ID] {
				ReceivedDuringInterval[netMessage.ID]++
				if ReceivedDuringInterval[netMessage.ID] < MessagesNeededWithinInterval {
					continue // not enough messages received yet
				} else {
					ReceivedIntervals[netMessage.ID]++
					ReceivedDuringInterval[netMessage.ID] = 0
					if ReceivedIntervals[netMessage.ID] < IntervalsNeeded {
						continue // not enough intervals
					} else {
						//Elevator is now back online
						NetError[netMessage.ID] = false
						netErrorToState <- NetErrorNotification{ID: netMessage.ID, NetError: false}
					}
				}
			}
			//If no neterror: process message normally
			resetTimer[netMessage.ID] <- struct{}{}

			//Avoid bothering state with duplicate messages
			if netMessage == prevNetMessages[netMessage.ID] {
				continue
			}

			netMessageToState <- netMessage
			prevNetMessages[netMessage.ID] = netMessage

		case id := <-timeout:
			//Handling a timout for an elevator that is online, involves marking as neterror
			if !NetError[id] {
				NetError[id] = true
				netErrorToState <- NetErrorNotification{ID: id, NetError: true}
			}
			//Timout during "get back online phase" means start from scratch.
			//Timout while online should also reset the timer.
			ReceivedDuringInterval[id] = 0
			ReceivedIntervals[id] = 0
			resetTimer[id] <- struct{}{}
		}
	}
}

// Monitors a timer for a specific elevator, and notifies
// the NetworkReceiver when the timer expires, and for which elevator it expired.
// The timer can be reset using the resetTimer channel.
func timeoutNotifier(id int, timeout chan int, resetTimer chan struct{}) {
	timer := time.NewTimer(time.Duration(NetErrorTimerLength * float64(time.Second)))
	for {
		select {
		case <-timer.C:
			timeout <- id
		case <-resetTimer:
			timer.Stop()
			timer.Reset(time.Duration(NetErrorTimerLength * float64(time.Second)))
		}
	}
}

// Reads raw data from the network, decodes it into netMessage structs,
// and forwards valid messages to NetworkReceiver.
// Filters out messeages from this elevator.
func receiver(receiveMessage chan<- NetMessage) {
	for {
		data, err := networkLow.Receive()
		if err != nil {
			continue
		}
		var netMessage NetMessage
		err = json.Unmarshal(data, &netMessage)
		if err != nil {
			log.Println("Receive json unmarshal error:", err)
			continue
		}
		if netMessage.ID == ID() {
			continue
		}

		receiveMessage <- netMessage
	}
}
