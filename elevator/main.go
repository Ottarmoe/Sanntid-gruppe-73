package main

import (
	// "elevator/network"
	// "elevator/networkLow"
	//"elevator/tests"
	. "elevator/elevatorConstants"
	referenceGenerator "elevator/refereceGenerator"
	"elevator/state"
	"elevio"
	//"elevator/hra"
)

func main() {
	// tests.TestMultipleServers()
	//tests.TimeHRA()
	//serverAdress := fmt.Sprintf("localhost:%d", 15657)
	elevio.Init("localhost:15657", NumFloors)

	sense_buttons := make(chan elevio.ButtonEvent)
	sense_floor := make(chan int)
	sense_obstr := make(chan bool)
	sense_stop := make(chan bool)
	int_mot := make(chan state.MotorState)
	int_mech := make(chan bool)

	stat_Gen := make(chan state.ElevWorldView)
	stat_Cont := make(chan state.ElevWorldView)
	stat_Insp := make(chan state.ElevWorldView)

	go elevio.PollButtons(sense_buttons)
	go elevio.PollFloorSensor(sense_floor)
	go elevio.PollObstructionSwitch(sense_obstr)
	go elevio.PollStopButton(sense_stop)

	go state.StateKeeper(0, 0, sense_buttons, sense_floor, int_mot, int_mech, stat_Gen, stat_Cont, stat_Insp)
	go referenceGenerator.ReferenceGenerator(stat_Gen)

	select {}
}
