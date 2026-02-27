package state

import (
	. "elevator/elevatorConstants"
	. "elevio"
	"fmt"
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
	Down Direction = iota
	Up
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
	motor MotorState
	floor int
}

type ElevState struct {
	NetError     bool
	CabAgreement [NumFloors]bool
	CabPriority  bool
	HallOrders   [NumFloors][2]hallOrderState //0 is down, 1 is up, use "direction"
	CabOrders    [NumFloors]cabOrderState
	CabPhysics   PhysicalState
	CabMechError bool
}

func PrintElevState(sta ElevState) {
	fmt.Printf(`
	NetError: %v
	CabAgreement: %v
	HallOrders: %v
	CabOrders: %v
	CabFloor: %v
	CabMotor: %v
	CabMechError: %v
	`, sta.NetError, sta.CabAgreement, sta.HallOrders, sta.CabOrders, sta.CabFloor, sta.CabMotor, sta.CabMechError)
}

type ElevWorldView struct {
	id    int
	elevs [NumElevators]ElevState
}

type NetMessage struct {
	HallOrders   [NumFloors][2]hallOrderState //0 is down, 1 is up, use "direction"
	CabOrders    [NumFloors]cabOrderState
	CabFloor     int
	CabMotor     MotorState
	CabMechError bool
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
) {

	var wView ElevWorldView
	wView.id = id
	for el := 0; el < NumElevators; el++ {
		wView.elevs[el].NetError = true
		wView.elevs[el].CabPriority = true

		for floor := 0; floor < NumFloors; floor++ {
			wView.elevs[el].HallOrders[floor][Down] = HallNO
			wView.elevs[el].HallOrders[floor][Up] = HallNO
			wView.elevs[el].CabOrders[floor] = CabUO
		}
	}
	//me.NetError = true //trust me bro

	me := &wView.elevs[id]
	me.CabMechError = false
	me.CabPhysics.motor.Behaviour = Idle
	me.CabPhysics.floor = initfloor

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
		updateLights(&consensus)

		stateComRefGenerator <- consensus
		//stateComController<-consensus
		//stateComInspector<-consensus
	}
}

func handleButton(vw *ElevWorldView, event ButtonEvent) {
	me := &vw.elevs[vw.id]
	switch event.Button {
	case BT_HallUp:
		if me.HallOrders[event.Floor][Up] == HallNO {
			me.HallOrders[event.Floor][Up] = HallO
		}
	case BT_HallDown:
		if me.HallOrders[event.Floor][Down] == HallNO {
			me.HallOrders[event.Floor][Down] = HallO
		}
	case BT_Cab:
		me.CabOrders[event.Floor] = CabO
		for _, elev := range vw.elevs {
			elev.CabAgreement[event.Floor] = false
		}
	}
}

func handleFloor(wv *ElevWorldView, event int) {
	wv.elevs[wv.id].CabPhysics.floor = event
}
func handleMotor(wv *ElevWorldView, event MotorState) {
	wv.elevs[wv.id].CabPhysics.motor.Behaviour = event.Behaviour
	if event.Behaviour == Moving {
		wv.elevs[wv.id].CabPhysics.motor.MovDirection = event.MovDirection
	}
}

func handleMech(wv *ElevWorldView, event bool) {
	wv.elevs[wv.id].CabMechError = event
}

func findConsensus(wv *ElevWorldView) ElevWorldView {
	consensus := *wv
	for floor := 0; floor < NumFloors; floor++ {
		hallDownExists := false
		hallUpExists := false
		cabExists := false
		elevExists := false
		for elev := 0; elev < NumElevators; elev++ {
			if wv.elevs[elev].NetError == false {
				elevExists = true
				if wv.elevs[elev].HallOrders[floor][Up] == HallO {
					hallDownExists = true
				}
				if wv.elevs[elev].HallOrders[floor][Down] == HallO {
					hallUpExists = true
				}
				if wv.elevs[elev].CabAgreement[floor] == true {
					cabExists = true
				}
			}
		}
		if (!elevExists || hallDownExists) && (wv.elevs[wv.id].HallOrders[floor][Down] == HallO) {
			consensus.elevs[consensus.id].HallOrders[floor][Down] = HallO
		} else {
			consensus.elevs[consensus.id].HallOrders[floor][Down] = HallNO
		}
		if (!elevExists || hallUpExists) && (wv.elevs[wv.id].HallOrders[floor][Up] == HallO) {
			consensus.elevs[consensus.id].HallOrders[floor][Up] = HallO
		} else {
			consensus.elevs[consensus.id].HallOrders[floor][Up] = HallNO
		}
		if (!elevExists || cabExists) && (wv.elevs[wv.id].CabOrders[floor] == CabO) {
			consensus.elevs[consensus.id].CabOrders[floor] = CabO
		} else {
			consensus.elevs[consensus.id].CabOrders[floor] = CabNO
		}
	}
	return consensus
}

func updateLights(consensus *ElevWorldView) {
	me := &consensus.elevs[consensus.id]
	for floor := 0; floor < NumFloors; floor++ {
		if me.HallOrders[floor][Down] == HallO {
			SetButtonLamp(BT_HallDown, floor, true)
		} else {
			SetButtonLamp(BT_HallDown, floor, false)
		}
		if me.HallOrders[floor][Up] == HallO {
			SetButtonLamp(BT_HallUp, floor, true)
		} else {
			SetButtonLamp(BT_HallUp, floor, false)
		}
		if me.CabOrders[floor] == CabO {
			SetButtonLamp(BT_Cab, floor, true)
		} else {
			SetButtonLamp(BT_Cab, floor, false)
		}
	}
	SetFloorIndicator(me.CabPhysics.floor)
}
