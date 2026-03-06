package hardwareControl

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
	. "elevio"
)

func HardWareControl(physicsToHardware <-chan PhysicalState, ordersWithConsesusToHardware <-chan OrdersWithConsesus) {
	for {
		select {
		case physicalState := <-physicsToHardware:
			SetFloorIndicator(physicalState.Floor)

			if physicalState.Behaviour == Idle {
				SetMotorDirection(0)
				// SetDoorOpenLamp(false)
			}

			if physicalState.Behaviour == Moving {
				// SetDoorOpenLamp(false)
				if physicalState.MovDirection == Up {
					SetMotorDirection(1)
				}
				if physicalState.MovDirection == Down {
					SetMotorDirection(-1)
				}
			}
			if physicalState.Behaviour == DoorOpen {
				SetMotorDirection(0)
				SetDoorOpenLamp(true)
			}

		case ordersWithConsesus := <-ordersWithConsesusToHardware:
			for floor := 0; floor < NumFloors; floor++ {
				SetButtonLamp(BT_HallDown, floor, ordersWithConsesus.HallOrders[floor][Down])
				SetButtonLamp(BT_HallUp, floor, ordersWithConsesus.HallOrders[floor][Up])
				SetButtonLamp(BT_Cab, floor, ordersWithConsesus.CabOrders[ordersWithConsesus.ID][floor])
			}
		}
	}
}
