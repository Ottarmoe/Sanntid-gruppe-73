package hardwareController

import (
	"elevator/elevatorConstants"
	. "elevator/stateTypes"
	. "elevio"
)

func HardwareController(physicalToHardware <-chan PhysicalState, ordersWithConsesusToHardware <-chan OrdersWithConsensus) {
	var prevConsensus OrdersWithConsensus
	for {
		select {
		case physicalState := <-physicalToHardware:
			physicalStateToHardware(physicalState)
		case ordersWithConsesus := <-ordersWithConsesusToHardware:
			ordersWithConsensusToHardware(ordersWithConsesus, prevConsensus)
			prevConsensus = ordersWithConsesus
		}
	}
}

// physicalStateToHardware drives motor and door lamp to match the current behaviour
func physicalStateToHardware(state PhysicalState) {
	SetFloorIndicator(state.Floor)
	switch state.Behaviour {
	case Idle:
		SetMotorDirection(MD_Stop)
		SetDoorOpenLamp(false)
	case Moving:
		SetDoorOpenLamp(false)
		if state.MovDirection == Up {
			SetMotorDirection(MD_Up)
		}
		if state.MovDirection == Down {
			SetMotorDirection(MD_Down)
		}
	case DoorOpen:
		SetMotorDirection(MD_Stop)
		SetDoorOpenLamp(true)
	}
}

// ordersWithConsensusToHardware updates button lamps to reflect agreed-upon orders.
// Only updates lamps that have changed since the last update to minimize redundant hardware calls.
func ordersWithConsensusToHardware(orders OrdersWithConsensus, prev OrdersWithConsensus) {
	for floor := 0; floor < elevatorConstants.NumFloors; floor++ {
		if orders.HallOrders[floor][Down] != prev.HallOrders[floor][Down] {
			SetButtonLamp(BT_HallDown, floor, orders.HallOrders[floor][Down])
		}
		if orders.HallOrders[floor][Up] != prev.HallOrders[floor][Up] {
			SetButtonLamp(BT_HallUp, floor, orders.HallOrders[floor][Up])
		}
		if orders.CabOrders[orders.ID][floor] != prev.CabOrders[orders.ID][floor] {
			SetButtonLamp(BT_Cab, floor, orders.CabOrders[orders.ID][floor])
		}
	}
}
