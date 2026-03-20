package main

import (
	. "elevator/elevatorConstants"
	"elevator/hardware"
	. "elevator/network"
	"elevator/networkLow"

	"elevator/hardwareLow"
	. "elevator/hardwareLow"
	"elevator/logicalControl"
	. "elevator/sharedTypes"
	"elevator/state"
	"time"
)

func main() {
	ConstantsInit() //always call this first to set up the constants, especially the elevatorID

	hardwareLow.Init()

	networkLow.Init()

	senseButtonsCh := make(chan ButtonEvent)
	senseFloorCh := make(chan int)
	senseObstructionCh := make(chan bool)

	motorStateCh := make(chan PhysicalState, 10)
	mechErrorCh := make(chan bool)

	refRequestCh := make(chan struct{}, 20)
	refToControllerCh := make(chan PhysicalState)
	statToControllerCh := make(chan PhysicalState, 10)

	netMessageToNetworkSenderCh := make(chan NetMessage, 10)
	netMessageToStateCh := make(chan NetMessage)
	netErrorToStateCh := make(chan NetErrorNotification)
	stateLifeSignalCh := make(chan struct{})

	ordersWithConsensusToHardwareCh := make(chan OrdersWithConsensus)
	physicsToHardwareCh := make(chan PhysicalState)

	startfloor := PhysicalInit()

	go hardware.PollButtons(senseButtonsCh)
	go hardware.PollFloorSensor(senseFloorCh)
	go hardware.PollObstructionSwitch(senseObstructionCh)

	go state.StateKeeper(startfloor,
		senseButtonsCh, senseFloorCh, motorStateCh, mechErrorCh,
		ordersWithConsensusToHardwareCh, physicsToHardwareCh,
		statToControllerCh, refRequestCh, refToControllerCh,
		netMessageToNetworkSenderCh, netMessageToStateCh, netErrorToStateCh,
		stateLifeSignalCh)
	go hardware.HardwareOut(physicsToHardwareCh, ordersWithConsensusToHardwareCh)
	go logicalControl.LogicalController(refToControllerCh, statToControllerCh, senseObstructionCh, refRequestCh, motorStateCh, mechErrorCh)
	go NetworkSender(netMessageToNetworkSenderCh)
	go NetworkReceiver(netMessageToStateCh, netErrorToStateCh)

	suicideWatchDog(stateLifeSignalCh)
}

func PhysicalInit() int {
	if GetFloor() != -1 {
		return GetFloor()
	}
	SetMotorDirection(MotorDirDown)
	for GetFloor() == -1 {
	}

	SetMotorDirection(MotorDirStop)
	return GetFloor()
}

func suicideWatchDog(stillAliveCh <-chan struct{}) {
	deathTimer := time.NewTimer(DeathCountDown)
	for {
		select {
		case <-deathTimer.C:
			panic("state timed out")
		case <-stillAliveCh:
			deathTimer = time.NewTimer(DeathCountDown)
		}
	}
}
