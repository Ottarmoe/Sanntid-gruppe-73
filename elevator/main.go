package main

import (
	// "elevator/network"
	// "elevator/networkLow"
	//"elevator/tests"
	. "elevator/elevatorConstants"
	"elevator/state"
	"elevio"
)

func main() {
	// tests.TestMultipleServers()
	//tests.TimeHRA()
	//serverAdress := fmt.Sprintf("localhost:%d", 15657)
	elevio.Init("localhost:15657", NumFloors)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	drv_mot := make(chan state.MotorState)
	drv_mech := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	state.FiniteStateMachine(0, 0, drv_buttons, drv_floors, drv_mot, drv_mech)

	select {}
}
