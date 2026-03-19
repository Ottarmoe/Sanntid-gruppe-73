package main

import (
	. "elevator/elevatorConstants"
	"elevator/hardwareControl"
	. "elevator/network"
	"elevator/networkLow"

	"elevator/hardwareLow"
	. "elevator/hardwareLow"
	"elevator/logicalControl"
	"elevator/state"
	. "elevator/stateTypes"
	"fmt"
)

func main() {
	ConstantsInit()

	serverAdress := fmt.Sprintf("localhost:%d", 15657+ID())
	hardwareLow.Init(serverAdress)
	networkLow.Init()

	sense_buttons := make(chan ButtonEvent)
	sense_floor := make(chan int)
	sense_obstr := make(chan bool)
	sense_stop := make(chan bool)
	int_mot := make(chan PhysicalState, 10)
	int_mech := make(chan bool)

	ref_request := make(chan struct{}, 20)
	ref_to_controller := make(chan PhysicalState)
	stat_to_controller := make(chan PhysicalState, 10)

	netMessageToNetworkSender := make(chan NetMessage, 10)
	netMessageToState := make(chan NetMessage)
	netErrorToState := make(chan NetErrorNotification)
	pokeStateCh := make(chan struct{})

	ordersWithConsensusToHardware := make(chan OrdersWithConsensus)
	physicsToHardware := make(chan PhysicalState)

	startfloor := PhysicalInit()

	go hardwareControl.PollButtons(sense_buttons)
	go hardwareControl.PollFloorSensor(sense_floor)
	go hardwareControl.PollObstructionSwitch(sense_obstr)
	go hardwareControl.PollStopButton(sense_stop)

	go state.StateKeeper(startfloor,
		sense_buttons, sense_floor, int_mot, int_mech,
		ordersWithConsensusToHardware, physicsToHardware,
		stat_to_controller, ref_request, ref_to_controller,
		netMessageToNetworkSender, netMessageToState, netErrorToState,
		pokeStateCh)
	go hardwareControl.HardWareControl(physicsToHardware, ordersWithConsensusToHardware)
	go logicalControl.LogicalController(ref_to_controller, stat_to_controller, sense_obstr, ref_request, int_mot, int_mech)
	go NetworkSender(netMessageToNetworkSender, pokeStateCh)
	go NetworkReceiver(netMessageToState, netErrorToState)

	select {}
}

func PhysicalInit() int {
	if GetFloor() != -1 {
		return GetFloor()
	}
	SetMotorDirection(MD_Down)
	for GetFloor() == -1 {
	}

	SetMotorDirection(MD_Stop)
	return GetFloor()
}
