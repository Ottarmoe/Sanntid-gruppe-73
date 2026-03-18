package hardwareControl

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
	. "elevio"
)

func HardWareControl(physicalToHardware <-chan PhysicalState, ordersWithConsensusToHardwarech <-chan OrdersWithConsensus) {
	var prevConsensus OrdersWithConsensus
	resetLights()
	for {
		select {
		case physicalState := <-physicalToHardware:
			physicalStateToHardware(physicalState)

		case ordersWithConsensus := <-ordersWithConsensusToHardwarech:
			ordersWithConsensusToHardware(ordersWithConsensus, prevConsensus)
			prevConsensus = ordersWithConsensus
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
// Only updates lamps that have changed since the last update to minimize redundant hardware call
func ordersWithConsensusToHardware(orders OrdersWithConsensus, prev OrdersWithConsensus) {
	for floor := 0; floor < NumFloors; floor++ {
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

func resetLights() {
	for floor := 0; floor < NumFloors; floor++ {
		SetButtonLamp(BT_HallDown, floor, false)
		SetButtonLamp(BT_HallUp, floor, false)
		SetButtonLamp(BT_Cab, floor, false)
	}
}
