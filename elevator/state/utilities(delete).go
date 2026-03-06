package state

import (
	// . "elevator/elevatorConstants"
	// "elevio"
	"fmt"
	. "elevator/stateTypes"
)

func PrintElevState(sta ElevState) {
	fmt.Printf(`
	OrderState:
	  HallOrders: %v
	  CabOrders: %v
	PhysicalState:
	  Behaviour: %v
	  MovDirection: %v
	  Floor: %v
	  MechError: %v
	`,
		sta.OrderState.HallOrders,
		sta.OrderState.CabOrders,
		sta.PhysicalState.Behaviour,
		sta.PhysicalState.MovDirection,
		sta.PhysicalState.Floor,
		sta.PhysicalState.MechError,
	)
}

// type NetMessage struct {
// 	HallOrders   [NumFloors][2]hallOrderState //0 is down, 1 is up, use "direction"
// 	CabOrders    [NumFloors]cabOrderState
// 	CabFloor     int
// 	CabMotor     MotorState
// 	CabMechError bool
// }
