package state

import (
	//. "elevator/elevatorConstants"
	. "elevio"
	"fmt"
)

type hallOrderState int

const (
	HallO hallOrderState = iota
	HallNO
	HallOPR
)

type cabOrderState int

const (
	CabO cabOrderState = iota
	CabNO
	CabUO //unknown order
)

type Direction int

const (
	Down Direction = iota
	Up
)

type MotorBehaviour int

const (
	Moving MotorBehaviour = iota
	Idle
	DoorOpen
)

type MotorState struct {
	Behaviour    MotorBehaviour
	MovDirection Direction
}

type ElevState struct {
	NetError     bool
	CabAgreement [4]bool
	HallOrders   [4][2]hallOrderState //0 is down, 1 is up, use "direction"
	CabOrders    [4]cabOrderState
	CabFloor     int
	CabMotor     MotorState
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
	elevs [3]ElevState
}

type NetMessage struct {
	HallOrders   [4][2]hallOrderState //0 is down, 1 is up, use "direction"
	CabOrders    [4]cabOrderState
	CabFloor     int
	CabMotor     MotorState
	CabMechError bool
}

// obstruction is not considered a state, and is handled internally by the door system
func FiniteStateMachine(
	id int,
	initfloor int,
	buttonClick <-chan ButtonEvent,
	floorReached <-chan int,
	motor <-chan MotorState,
	mechError <-chan bool,
) {

	var wView ElevWorldView
	wView.id = id
	me := &wView.elevs[id]
	me.NetError = true //trust me bro

	me.CabMechError = false
	me.CabMotor.Behaviour = Idle
	me.CabFloor = initfloor

	for floor := 0; floor < 4; floor++ {
		me.HallOrders[floor][Down] = HallNO
		me.HallOrders[floor][Up] = HallNO
		me.CabOrders[floor] = CabUO
	}

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
		// case a := <- obstructionChange:
		// 	fmt.Printf("%+v\n", a)

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

func handleFloor(vw *ElevWorldView, event int) {
	vw.elevs[vw.id].CabFloor = event
}
func handleMotor(vw *ElevWorldView, event MotorState) {
	vw.elevs[vw.id].CabMotor.Behaviour = event.Behaviour
	if event.Behaviour == Moving {
		vw.elevs[vw.id].CabMotor.MovDirection = event.MovDirection
	}
}

func handleMech(vw *ElevWorldView, event bool) {
	vw.elevs[vw.id].CabMechError = event
}
