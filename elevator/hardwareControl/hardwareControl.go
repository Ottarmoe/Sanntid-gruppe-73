package hardwareControl

import (
	. "elevator/elevatorConstants"
	. "elevio"
	. "elevator/state"
)

func HardWareControl(physicsToHardware <-chan PhysicalState, ordersWithConsesusToHardware <-chan OrdersWithConsesus) {
	for {
		select {
		case physicalState := <-physicsToHardware:
			SetFloorIndicator(physicalState.Floor)

		case ordersWithConsesus := <-ordersWithConsesusToHardware:
			for floor := 0; floor < NumFloors; floor++ {
				SetButtonLamp(BT_HallDown, floor, ordersWithConsesus.HallOrders[floor][Down])
				SetButtonLamp(BT_HallUp, floor, ordersWithConsesus.HallOrders[floor][Up])
				SetButtonLamp(BT_Cab, floor, ordersWithConsesus.HallOrders[floor][Down])
			}	
		}
	}
}