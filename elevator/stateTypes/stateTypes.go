package stateTypes

import (
	. "elevator/elevatorConstants"
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
	ID         int
	HallOrders [NumFloors][2]bool //0 is down, 1 is up, use "direction"
	CabOrders  [NumElevators][NumFloors]bool
}

type OurOrders struct {
	HallOrders [NumFloors][2]bool //0 is down, 1 is up, use "direction"
	CabOrders  [NumFloors]bool
}