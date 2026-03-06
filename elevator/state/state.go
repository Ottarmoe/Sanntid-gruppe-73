package state

import (
	. "elevator/elevatorConstants"
	. "elevio"
)

// Types
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

// Building the world view struct
type PhysicalState struct {
	Behaviour    MotorBehaviour
	MovDirection Direction
	Floor        int
	MechError    bool
}

type OrderState struct {
	HallOrders [NumFloors][2]hallOrderState //0 is down, 1 is up, use "direction"
	CabOrders  [NumFloors]cabOrderState
}

type ElevState struct { //States to be mirrored to other elevators
	OrderState    OrderState
	PhysicalState PhysicalState
}

type ElevWorldView struct {
	ID             int
	ElevStates     [NumElevators]ElevState
	CabArchiveSeen [NumElevators]bool
	CabAgreement   [NumElevators][NumFloors]bool
	NetError       [NumElevators]bool
}

// Consesus struct
type OrdersWithConsesus struct {
	ID int
	HallOrders [NumFloors][2]bool //0 is down, 1 is up, use "direction"
	CabOrders  [NumElevators][NumFloors]bool
}

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

	stateToController chan<- ElevWorldView,
	referenceRequest chan<- struct{},
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
