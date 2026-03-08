package state

import (
	// . "elevator/elevatorConstants"
	// "elevio"
	. "elevator/stateTypes"
	"fmt"
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

func PrintNetMessage(msg NetMessage) {
	fmt.Printf(`
	ID: %v
	ElevState:
  OrderState:
	  HallOrders: %v
	  CabOrders: %v
	PhysicalState:
	  Behaviour: %v
	  MovDirection: %v
	  Floor: %v
	  MechError: %v
	CabBackups: %v
	`,
		msg.ID,
		msg.ElevState.OrderState.HallOrders,
		msg.ElevState.OrderState.CabOrders,
		msg.ElevState.PhysicalState.Behaviour,
		msg.ElevState.PhysicalState.MovDirection,
		msg.ElevState.PhysicalState.Floor,
		msg.ElevState.PhysicalState.MechError,
		msg.CabBackups,
	)
}