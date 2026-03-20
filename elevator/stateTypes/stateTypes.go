package stateTypes

import (
	. "elevator/elevatorConstants"
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

// Structs that represent purely if an order exists or not. 
type OrdersWithConsensus struct {
	HallOrders [NumFloors][2]bool 
	CabOrders  [NumElevators][NumFloors]bool
}

type AssignedOrders struct {
	HallOrders [NumFloors][2]bool 
	CabOrders  [NumFloors]bool
}

// Hardware event types
type MotorDirection int

const (
	MotorDirUp   MotorDirection = 1
	MotorDirDown                = -1
	MotorDirStop                = 0
)

type ButtonType int

const (
	ButtonHallUp   ButtonType = 0
	ButtonHallDown            = 1
	ButtonCab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

// helper functions
func (wv *ElevWorldView) IsOnline(elev int) bool {
	return !wv.NetError[elev] || elev == MyID()
}
func (wv *ElevWorldView) IsOnlinePeer(elev int) bool {
	return !wv.NetError[elev] && elev != MyID()
}
func (wv *ElevWorldView) IsOfflinePeer(elev int) bool {
	return wv.NetError[elev] && elev != MyID()
}
func (wv *ElevWorldView) AnyPeerExists() bool {
	for elev := 0; elev < NumElevators; elev++ {
		if wv.IsOnlinePeer(elev) {
			return true
		}
	}
	return false
}

func (wv *ElevWorldView) GetHallOrder(elev int, floor int, dir Direction) HallOrderState {
	return wv.ElevStates[elev].OrderState.HallOrders[floor][dir]
}

func (wv *ElevWorldView) MyElev() *ElevState {
	return &wv.ElevStates[MyID()]
}

func (wv *ElevWorldView) AnyoneInHallOrderState(hallOrderState HallOrderState, floor int, dir Direction) bool {
	for elev := 0; elev < NumElevators; elev++ {
		if wv.IsOnline(elev) {
			if wv.GetHallOrder(elev, floor, dir) == hallOrderState {
				return true
			}
		}
	}
	return false
}

func (wv *ElevWorldView) AnyoneElseInHallOrderState(hallOrderState HallOrderState, floor int, dir Direction) bool {
	for elev := 0; elev < NumElevators; elev++ {
		if wv.IsOnlinePeer(elev) {
			if wv.GetHallOrder(elev, floor, dir) == hallOrderState {
				return true
			}
		}
	}
	return false
}

func (wv *ElevWorldView) CabOrderArchiveExists(floor int) bool {
	for elev := 0; elev < NumElevators; elev++ {
		if wv.IsOnlinePeer(elev) {
			if wv.CabAgreement[elev][floor] {
				return true
			}
		}
	}
	return false
}

func (wv *ElevWorldView) CompileNetMessage() NetMessage {
	var cabBackups [NumElevators][NumFloors]CabOrderState
	for elev := 0; elev < NumElevators; elev++ {
		cabBackups[elev] = wv.ElevStates[elev].OrderState.CabOrders
	}
	netMessage := NetMessage{
		ID:         MyID(),
		ElevState:  *wv.MyElev(),
		CabBackups: cabBackups,
	}
	return netMessage
}
