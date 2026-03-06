package state

import (
	. "elevator/elevatorConstants"
	. "elevio"
	. "elevator/stateTypes"
)

// obstruction is not considered a state, and is handled internally by the door system
func StateKeeper(
	id int,
	initfloor int,
	buttonClick <-chan ButtonEvent,
	floorReached <-chan int,
	motor <-chan PhysicalState,
	mechError <-chan bool,

	ordersWithConsesusToHardware chan<- OrdersWithConsesus,
	physicsToHardware chan<- PhysicalState,

	stateToController chan<- PhysicalState,
	referenceRequest chan<- struct{},
	refToController chan<- PhysicalState,
) {

	var wView ElevWorldView = initWorldView(id, initfloor)
	me := &wView.ElevStates[id]
	physicalState := &me.PhysicalState

	for {
		PrintElevState(*me)
		select {
		case buttonEvent := <-buttonClick:
			handleButton(&wView, buttonEvent)
		case floorEvent := <-floorReached:
			handleFloor(physicalState, floorEvent)
		case motorEvent := <-motor:
			handleMotor(&wView, motorEvent)
		case mechEvent := <-mechError:
			handleMech(&wView, mechEvent)
		}
		ordersWithConsesus := findConsensus(wView)

		ordersWithConsesusToHardware <- ordersWithConsesus
		physicsToHardware <- *physicalState

		// stateComRefGenerator <- consensus
		//stateComController<-consensus
		//stateComInspector<-consensus
	}
}

func initWorldView(id int, initfloor int) ElevWorldView {
	var wView ElevWorldView

	wView.ID = id
	for elev := 0; elev < NumElevators; elev++ {
		wView.NetError[elev] = true
		wView.CabArchiveSeen[elev] = false

		for floor := 0; floor < NumFloors; floor++ {
			wView.ElevStates[elev].OrderState.HallOrders[floor][Down] = HallNO
			wView.ElevStates[elev].OrderState.HallOrders[floor][Up] = HallNO
			wView.ElevStates[elev].OrderState.CabOrders[floor] = CabUO
		}
	}

	me := &wView.ElevStates[id]
	me.PhysicalState.MechError = false
	me.PhysicalState.Behaviour = Idle
	me.PhysicalState.Floor = initfloor

	return wView
}
