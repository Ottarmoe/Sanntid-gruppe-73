package hardwareControl

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
	. "elevio"
)

func HardWareControl(physicsToHardware <-chan PhysicalState, ordersWithConsesusToHardware <-chan OrdersWithConsesus) {
	var prevConsensus OrdersWithConsesus
	for {
		select {
		case physicalState := <-physicsToHardware:
			SetFloorIndicator(physicalState.Floor)

			if physicalState.Behaviour == Idle {
				SetMotorDirection(0)
				SetDoorOpenLamp(false)
			}

			if physicalState.Behaviour == Moving {
				SetDoorOpenLamp(false)
				//fmt.Println("moving", physicalState.MovDirection)
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
				if ordersWithConsesus.HallOrders[floor][Down] != prevConsensus.HallOrders[floor][Down] {
					SetButtonLamp(BT_HallDown, floor, ordersWithConsesus.HallOrders[floor][Down])
				}
				if ordersWithConsesus.HallOrders[floor][Up] != prevConsensus.HallOrders[floor][Up] {
					SetButtonLamp(BT_HallDown, floor, ordersWithConsesus.HallOrders[floor][Up])
				}
				if ordersWithConsesus.CabOrders[ordersWithConsesus.ID][floor] != prevConsensus.CabOrders[ordersWithConsesus.ID][floor] {
					SetButtonLamp(BT_Cab, floor, ordersWithConsesus.CabOrders[ordersWithConsesus.ID][floor])
				}
			}
			prevConsensus = ordersWithConsesus
		}
	}
}
