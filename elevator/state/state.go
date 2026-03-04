package state

import (
	. "elevator/elevatorConstants"
	. "elevio"
)

//Types
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

//Building the world view struct
type PhysicalState struct {
	Behaviour    MotorBehaviour
	MovDirection Direction
	Floor int
}

type OrderState struct {
	CabAgreement [NumFloors]bool
	CabPriority  bool
	HallOrders   [NumFloors][2]hallOrderState //0 is down, 1 is up, use "direction"
	CabOrders    [NumFloors]cabOrderState
}

type ElevState struct {
	OrderState    OrderState
	PhysicalState PhysicalState
	NetError      bool
	MechError     bool
}

type ElevWorldView struct {
	ID    int
	Elevs [NumElevators]ElevState
}

//Consesus struct
type OrdersWithConsesus struct {
	HallOrders   [NumFloors][2]bool //0 is down, 1 is up, use "direction"
	CabOrders    [NumFloors]bool
}


// obstruction is not considered a state, and is handled internally by the door system
func StateKeeper(
	id int,
	initfloor int,
	buttonClick <-chan ButtonEvent,
	floorReached <-chan int,
	motor <-chan PhysicalState,
	mechError <-chan bool,
	stateComRefGenerator chan<- ElevWorldView,
	stateComController chan<- ElevWorldView,
	stateComInspector chan<- ElevWorldView,
	ordersWithConsesusToHardware chan<- OrdersWithConsesus,
    physicsToHardware chan<- PhysicalState,
	) {

	var wView ElevWorldView = initWorldView(id,initfloor);
	elevator := &wView.Elevs[id]	
	physicalState := &elevator.PhysicalState

	for {
		PrintElevState(*elevator)
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
		ordersWithConsesus := findConsensus(&wView)

		ordersWithConsesusToHardware <- ordersWithConsesus
		physicsToHardware <- physics

		// stateComRefGenerator <- consensus
		//stateComController<-consensus
		//stateComInspector<-consensus
	}
}


func initWorldView(id int, initfloor int) ElevWorldView {
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

	return wView
}
