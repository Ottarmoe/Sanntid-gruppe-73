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
func NetworkSender(netMessageToNetworkSenderCh <-chan NetMessage) {
	timeToSend := time.NewTicker(BroadcastRate)

	//No periodic sending before the first netmessage has been comunnicated by state
	netMessage := <-netMessageToNetworkSenderCh

	for {
		select {
		case netMessage = <-netMessageToNetworkSenderCh:
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
func NetworkReceiver(netMessageToStateCh chan<- NetMessage, netErrorToStateCh chan<- NetErrorNotification) {
	var prevNetMessages [NumElevators]NetMessage

	receiveMessageCh := make(chan NetMessage)
	go receiver(receiveMessageCh)

	//Track network error state and received messages per interval, nr. of intervals
	var NetError [NumElevators]bool
	var ReceivedDuringInterval [NumElevators]int
	var ReceivedIntervals [NumElevators]int
	for i := 0; i < NumElevators; i++ {
		NetError[i] = true
		ReceivedDuringInterval[i] = 0
		ReceivedIntervals[i] = 0
	}

	// Start a timeoutCh notifier for each other elevator
	timeoutCh := make(chan int)
	var resetTimerCh [NumElevators]chan struct{}
	for i := 0; i < NumElevators; i++ {
		if i == MyID() {
			continue
		}
		resetTimerCh[i] = make(chan struct{}, 1)
		go timeoutNotifier(i, timeoutCh, resetTimerCh[i])
	}

	for {
		select {
		case netMessage := <-receiveMessageCh:
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
						netErrorToStateCh <- NetErrorNotification{ID: netMessage.ID, NetError: false}
					}
				}
			}
			//If no neterror: process message normally
			resetTimerCh[netMessage.ID] <- struct{}{}

			//Avoid bothering state with duplicate messages
			if netMessage == prevNetMessages[netMessage.ID] {
				continue
			}

			netMessageToStateCh <- netMessage
			prevNetMessages[netMessage.ID] = netMessage

		case id := <-timeoutCh:
			//Handling a timout for an elevator that is online, involves marking as neterror
			if !NetError[id] {
				NetError[id] = true
				netErrorToStateCh <- NetErrorNotification{ID: id, NetError: true}
			}
			//Timout during "get back online phase" means start from scratch.
			//Timout while online should also reset the timer.
			ReceivedDuringInterval[id] = 0
			ReceivedIntervals[id] = 0
			resetTimerCh[id] <- struct{}{}
		}
	}
}

// Monitors a timer for a specific elevator, and notifies
// the NetworkReceiver when the timer expires, and for which elevator it expired.
// The timer can be reset using the resetTimer channel.
func timeoutNotifier(id int, timeoutCh chan int, resetTimerCh chan struct{}) {
	timer := time.NewTimer(time.Duration(NetErrorTimerLength))
	for {
		select {
		case <-timer.C:
			timeoutCh <- id
		case <-resetTimerCh:
			timer.Stop()
			timer.Reset(time.Duration(NetErrorTimerLength))
		}
	}
}

// Reads raw data from the network, decodes it into netMessage structs,
// and forwards valid messages to NetworkReceiver.
// Filters out messeages from this elevator.
func receiver(receiveMessageCh chan<- NetMessage) {
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
		if netMessage.ID == MyID() {
			continue
		}

		receiveMessageCh <- netMessage
	}
}
