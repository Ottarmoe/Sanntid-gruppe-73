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
	timeToSend := make(chan struct{})
	go sendRateTimer(timeToSend)
	
	netMessage := <- netMessageToNetworkSender;

	for{
		select {
		case netMessage = <- netMessageToNetworkSender:
			
		case <- timeToSend:	
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

		if(netMessage == prevNetMessage){
			continue
		}

		netMessageToState <- netMessage;
		prevNetMessage = netMessage
	}
}

func sendRateTimer(timeToSend chan<- struct{}) {
	for {
		time.Sleep(time.Second)
		timeToSend <- struct{}{}
	}
}
