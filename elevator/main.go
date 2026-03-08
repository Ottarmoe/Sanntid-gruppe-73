package main

import (
	. "elevator/elevatorConstants"
	. "elevator/hardwareControl"
	. "elevator/network"
	"elevator/networkLow"

	// "elevator/referenceGenerator"
	"elevator/logicalController"
	"elevator/state"
	. "elevator/stateTypes"
	"elevio"
	"fmt"
)

func main() {
	ConstantsInit()
	serverAdress := fmt.Sprintf("localhost:%d", 15657+ID())
	elevio.Init(serverAdress, NumFloors)
	networkLow.Init()

	sense_buttons := make(chan elevio.ButtonEvent)
	sense_floor := make(chan int)
	sense_obstr := make(chan bool)
	sense_stop := make(chan bool)
	int_mot := make(chan PhysicalState, 10)
	int_mech := make(chan bool)

	ref_request := make(chan struct{})
	ref_to_controller := make(chan PhysicalState)
	stat_to_controller := make(chan PhysicalState, 10)

	netMessageToNetworkSender := make(chan NetMessage)
	netMessageToState := make(chan NetMessage)

	ordersWithConsesusToHardware := make(chan OrdersWithConsesus)
	physicsToHardware := make(chan PhysicalState)

	go elevio.PollButtons(sense_buttons)
	go elevio.PollFloorSensor(sense_floor)
	go elevio.PollObstructionSwitch(sense_obstr)
	go elevio.PollStopButton(sense_stop)

	go state.StateKeeper(ID(), 0,
		sense_buttons, sense_floor, int_mot, int_mech,
		ordersWithConsesusToHardware, physicsToHardware,
		stat_to_controller, ref_request, ref_to_controller,
		netMessageToNetworkSender,netMessageToState)
	go HardWareControl(physicsToHardware, ordersWithConsesusToHardware)
	go logicalController.Controller(ref_to_controller, stat_to_controller, sense_obstr, ref_request, int_mot, int_mech)
	go NetworkSender(netMessageToNetworkSender)
	go NetworkReceiver(netMessageToState)
	
	// go referenceGenerator.ReferenceGenerator(stat_Gen)

	select {}
}
