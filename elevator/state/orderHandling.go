package state

import (
	. "elevator/elevatorConstants"
	. "elevio"
	// "fmt"
)

func handleButton(vw *ElevWorldView, event ButtonEvent) {
	me := &vw.ElevStates[vw.ID]
	switch event.Button {
	case BT_HallUp:
		if me.OrderState.HallOrders[event.Floor][Up] == HallNO {
			me.OrderState.HallOrders[event.Floor][Up] = HallO
		}
	case BT_HallDown:
		if me.OrderState.HallOrders[event.Floor][Down] == HallNO {
			me.OrderState.HallOrders[event.Floor][Down] = HallO
		}
	case BT_Cab:
		me.OrderState.CabOrders[event.Floor] = CabO
		for elev := 0; elev < NumElevators; elev++ {
			vw.CabAgreement[elev][event.Floor] = false
		}
	}
}

func findConsensus(wv ElevWorldView) OrdersWithConsesus {
	var ordersWithConsesus OrdersWithConsesus


	for floor := 0; floor < NumFloors; floor++ {
		hallDownExists := false
		hallUpExists := false
		cabExists := false
		elevExists := false
		for elev := 0; elev < NumElevators; elev++ {
			peerHallOrders := &wv.ElevStates[elev].OrderState.HallOrders

			if wv.NetError[elev] == false {
				elevExists = true
				if peerHallOrders[floor][Up] == HallO {
					hallDownExists = true
				}
				if peerHallOrders[floor][Down] == HallO {
					hallUpExists = true
				}
				if wv.CabAgreement[elev][floor] == true {
					cabExists = true
				}
			}
		}
		HallOrders := &wv.ElevStates[wv.ID].OrderState.HallOrders
		CabOrders := &wv.ElevStates[wv.ID].OrderState.CabOrders

		if (!elevExists || hallDownExists) && (HallOrders[floor][Down] == HallO) {
			ordersWithConsesus.HallOrders[floor][Down] = true
		} else {
			ordersWithConsesus.HallOrders[floor][Down] = false
		}
		if (!elevExists || hallUpExists) && (HallOrders[floor][Up] == HallO) {
			ordersWithConsesus.HallOrders[floor][Up] = true
		} else {
			ordersWithConsesus.HallOrders[floor][Up] = false
		}
		if (!elevExists || cabExists) && (CabOrders[floor] == CabO) {
			ordersWithConsesus.CabOrders[floor] = true
		} else {
			ordersWithConsesus.CabOrders[floor] = false
		}
	}
	return ordersWithConsesus
}