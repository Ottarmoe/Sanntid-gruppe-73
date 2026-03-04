package hardwareControl

import (
	. "elevator/elevatorConstants"
	. "elevio"
	. "elevator/state"
)

func HardWareControl(wv <-chan ElevWorldView) {
	for {
		consensus := <- wv;

		me := &consensus.Elevs[consensus.ID]
		for floor := 0; floor < NumFloors; floor++ {
			if me.HallOrders[floor][Down] == HallO {
				SetButtonLamp(BT_HallDown, floor, true)
			} else {
				SetButtonLamp(BT_HallDown, floor, false)
			}
			if me.HallOrders[floor][Up] == HallO {
				SetButtonLamp(BT_HallUp, floor, true)
			} else {
				SetButtonLamp(BT_HallUp, floor, false)
			}
			if me.CabOrders[floor] == CabO {
				SetButtonLamp(BT_Cab, floor, true)
			} else {
				SetButtonLamp(BT_Cab, floor, false)
			}
		}
		SetFloorIndicator(me.CabPhysics.Floor)
	}

}