package state

import (
	. "elevator/elevatorConstants"
	. "elevio"
)

type hallOrderState int

const (
	HallNO hallOrderState = iota
	HallO
	HallOPR
)

type cabOrderState int

const (
	CabNO cabOrderState = iota
	CabO
	CabUO //unknown order
)

type Direction int

const (
	Up Direction = iota
	Down
)

type MotorBehaviour int

const (
	Idle MotorBehaviour = iota
	Moving
	DoorOpen
)

type MotorState struct {
	Behaviour    MotorBehaviour
	MovDirection Direction
}

type PhysicalState struct {
	Motor MotorState
	Floor int
}

type ElevState struct {
	//Order relevant
	NetError     bool
	CabAgreement [NumFloors]bool
	CabPriority  bool
	HallOrders   [NumFloors][2]hallOrderState //0 is down, 1 is up, use "direction"
	CabOrders    [NumFloors]cabOrderState
	//Physics
	CabPhysics   PhysicalState
	CabMechError bool
}

type ElevWorldView struct {
	ID    int
	Elevs [NumElevators]ElevState
}

// obstruction is not considered a state, and is handled internally by the door system
func StateKeeper(
	id int,
	initfloor int,
	buttonClick <-chan ButtonEvent,
	floorReached <-chan int,
	motor <-chan MotorState,
	mechError <-chan bool,
	stateComRefGenerator chan<- ElevWorldView,
	stateComController chan<- ElevWorldView,
	stateComInspector chan<- ElevWorldView,
	hardWareControl chan<- ElevWorldView,
	) {

	var wView ElevWorldView
	wView.ID = id
	for elev := 0; elev < NumElevators; elev++ {
		wView.Elevs[elev].NetError = true
		wView.Elevs[elev].CabPriority = true

		for floor := 0; floor < NumFloors; floor++ {
			wView.Elevs[elev].HallOrders[floor][Down] = HallNO
			wView.Elevs[elev].HallOrders[floor][Up] = HallNO
			wView.Elevs[elev].CabOrders[floor] = CabUO
		}
	}

	me := &wView.Elevs[id]
	me.NetError = true //trust me bro
	me.CabMechError = false
	me.CabPhysics.Motor.Behaviour = Idle
	me.CabPhysics.Floor = initfloor

	for {
		PrintElevState(*me)
		select {
		case buttonEvent := <-buttonClick:
			handleButton(&wView, buttonEvent)
		case floorEvent := <-floorReached:
			handleFloor(&wView, floorEvent)
		case motorEvent := <-motor:
			handleMotor(&wView, motorEvent)
		case mechEvent := <-mechError:
			handleMech(&wView, mechEvent)
		}
		consensus := findConsensus(&wView)
		hardWareControl <- consensus

		stateComRefGenerator <- consensus
		//stateComController<-consensus
		//stateComInspector<-consensus
	}
}


