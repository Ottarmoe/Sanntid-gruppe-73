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

type OrderState struct {
	HallOrders [NumFloors][2]HallOrderState //0 is up, 1 is down, use "direction"
	CabOrders  [NumFloors]CabOrderState
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

// Consesus struct
type OrdersWithConsesus struct {
	ID         int
	HallOrders [NumFloors][2]bool //0 is up, 1 is down, use "direction"
	CabOrders  [NumElevators][NumFloors]bool
}

type OurOrders struct {
	HallOrders [NumFloors][2]bool //0 is up, 1 is down, use "direction"
	CabOrders  [NumFloors]bool
}
