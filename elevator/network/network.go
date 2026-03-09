package network

import (
	"elevator/networkLow"
	// "fmt"
	. "elevator/stateTypes"
	"time"

	// "elevator/elevatorConstants"
	. "elevator/elevatorConstants"
	"encoding/json"
	"log"
)

func NetworkSender(netMessageToNetworkSender <-chan NetMessage){
	timeToSend := time.NewTicker(1 * time.Second)
	
	netMessage := <- netMessageToNetworkSender;

	for{
		select {
		case netMessage = <- netMessageToNetworkSender:
			
		case <- timeToSend.C:	
			    data, err := json.Marshal(netMessage)
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

func NetworkReceiver(netMessageToState chan<- NetMessage){
	var prevNetMessage NetMessage

	receiveMessage := make(chan NetMessage)
	go receiver(receiveMessage)

	var NetError [NumElevators]bool
	for i := range NumElevators {
    NetError[i] = true
	}
	var timers [NumElevators]*time.Timer
	for i := 0; i < NumElevators; i++ {
		if(i == ID()){
			continue
		}
		timers[i] = time.NewTimer(5 * time.Second)
	}

	for{
	select{
		case netMessage := <- receiveMessage:
			//Handle neterror
			if(NetError[netMessage.ID]){
				timers[netMessage.ID].Reset(1 * time.Second)

			}
			timers[netMessage.ID].Reset(5 * time.Second)

			//Avoid bothering state with duplicate messages
			if(netMessage == prevNetMessage){
				continue
			}

			netMessageToState <- netMessage;
			prevNetMessage = netMessage
		case 
		
	}
	}
}

func receiver(receiveMessage chan<- NetMessage) {
	for{
		data, err := networkLow.Receive()
		if err != nil {
			log.Println("receive  error:", err) 
			continue
		}
		var netMessage NetMessage
		err = json.Unmarshal(data, &netMessage)
		if err != nil {
			log.Println("receive json unmarshal error:", err) 
			continue
		}
		if(netMessage.ID == ID()){
			continue
		}
		receiveMessage <- netMessage
	}

}