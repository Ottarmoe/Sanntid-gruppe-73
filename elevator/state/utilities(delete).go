package state

import (
	// . "elevator/elevatorConstants"
	// "elevio"
	"fmt"
)

func PrintElevState(sta ElevState) {
	fmt.Printf(`
	NetError: %v
	CabAgreement: %v
	HallOrders: %v
	CabOrders: %v
	CabFloor: %v
	CabMotor: %v
	CabMechError: %v
	`, sta.NetError, sta.CabAgreement, sta.HallOrders, sta.CabOrders, sta.CabPhysics.Floor, sta.CabPhysics.Motor, sta.CabMechError)
}

// type NetMessage struct {
// 	HallOrders   [NumFloors][2]hallOrderState //0 is down, 1 is up, use "direction"
// 	CabOrders    [NumFloors]cabOrderState
// 	CabFloor     int
// 	CabMotor     MotorState
// 	CabMechError bool
// }