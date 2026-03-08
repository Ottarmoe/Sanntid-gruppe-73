package network

import (
	"elevator/networkLow"
	// "fmt"
	"time"
	. "elevator/stateTypes"
	// "elevator/elevatorConstants"
	"encoding/json"
	"log"
	. "elevator/elevatorConstants"
)

func NetworkCommunicator(netMessageToNetworkCommunicator <-chan NetMessage){
	timeToSend := make(chan struct{})
	go sendRateTimer(timeToSend)
	
	netMessage := <- netMessageToNetworkCommunicator;

	for{
		select {
		case netMessage = <- netMessageToNetworkCommunicator:
			
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

		netMessageToState <- netMessage;
	}
}

func sendRateTimer(timeToSend chan<- struct{}) {
	for {
		time.Sleep(time.Second)
		timeToSend <- struct{}{}
	}
}
