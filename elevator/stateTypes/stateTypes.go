package stateTypes

import (
	. "elevator/elevatorConstants"
	"fmt"
)

// Types
type HallOrderState int

const (
	HallNO HallOrderState = iota
	HallO
	HallOPR
)

type CabOrderState int

const (
	CabNO CabOrderState = iota
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
	HallOrders [NumFloors][2]HallOrderState //0 is up, 1 is down, use "direction"
	CabOrders  [NumFloors]CabOrderState
}

type ElevState struct { //States to be mirrored to other elevators
	OrderState    OrderState
	PhysicalState PhysicalState
}

type ElevWorldView struct {
	ElevStates     [NumElevators]ElevState
	CabArchiveSeen [NumElevators]bool
	CabAgreement   [NumElevators][NumFloors]bool
	NetError       [NumElevators]bool
}

// Network
type NetMessage struct {
	ID         int
	ElevState  ElevState
	CabBackups [NumElevators][NumFloors]CabOrderState //not needed by other elevators most of the time
}

type NetErrorNotification struct {
	ID       int
	NetError bool
}

// Consensus struct
type OrdersWithConsensus struct {
	HallOrders [NumFloors][2]bool //0 is up, 1 is down, use "direction"
	CabOrders  [NumElevators][NumFloors]bool
}

type OurOrders struct {
	HallOrders [NumFloors][2]bool //0 is up, 1 is down, use "direction"
	CabOrders  [NumFloors]bool
}

func PrintPhysicalState(stat PhysicalState) {
	switch stat.Behaviour {
	case Idle:
		fmt.Print("Idle ", []string{"Up", "Down"}[stat.MovDirection])
	case Moving:
		fmt.Print("Moving ", []string{"Up", "Down"}[stat.MovDirection])
	case DoorOpen:
		fmt.Print("DoorOpen ", []string{"Up", "Down"}[stat.MovDirection])
	}
	fmt.Println(" on floor", stat.Floor)
}

//Hardware event types
type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

//helper functions
func (wv *ElevWorldView) IsOnline(elev int) bool {
    return !wv.NetError[elev] || elev == ID()
}

func (wv *ElevWorldView) GetHallOrder(elev int, floor int, dir Direction) HallOrderState {
    return wv.ElevStates[elev].OrderState.HallOrders[floor][dir]
}

func (wv *ElevWorldView) MyState() *ElevState {
    return &wv.ElevStates[ID()]
}

func (wv *ElevWorldView) AnyoneInHallOPR(floor int, dir Direction) bool {
    for elev := 0; elev < NumElevators; elev++ {
        if wv.IsOnline(elev) {
            if wv.GetHallOrder(elev, floor, dir) == HallOPR {
                return true
            }
        }
    }
    return false
}