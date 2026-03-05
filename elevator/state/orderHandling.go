package state

import (
	. "elevator/elevatorConstants"
	. "elevio"
	// "fmt"
)

func handleButton(vw *ElevWorldView, event ButtonEvent) {
	me := &vw.Elevs[vw.ID]
	switch event.Button {
	case BT_HallUp:
		if me.HallOrders[event.Floor][Up] == HallNO {
			me.HallOrders[event.Floor][Up] = HallO
		}
	case BT_HallDown:
		if me.HallOrders[event.Floor][Down] == HallNO {
			me.HallOrders[event.Floor][Down] = HallO
		}
	case BT_Cab:
		me.CabOrders[event.Floor] = CabO
		for _, elev := range vw.Elevs {
			elev.CabAgreement[event.Floor] = false
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
			peerCabOrders := &wv.ElevStates[elev].OrderState.CabOrders


			if wv.Elevs[elev].NetError == false {
				elevExists = true
				if peerHallOrders[floor][Up] == HallO {
					hallDownExists = true
				}
				if peerHallOrders[floor][Down] == HallO {
					hallUpExists = true
				}
				if wv.Elevs[elev].CabAgreement[floor] == true {
					cabExists = true
				}
			}
		}
		if (!elevExists || hallDownExists) && (wv.Elevs[wv.ID].HallOrders[floor][Down] == HallO) {
			ordersWithConsesus.HallOrders[floor][Down] = true
		} else {
			ordersWithConsesus.HallOrders[floor][Down] = false
		}
		if (!elevExists || hallUpExists) && (wv.Elevs[wv.ID].HallOrders[floor][Up] == HallO) {
			ordersWithConsesus.HallOrders[floor][Up] = true
		} else {
			ordersWithConsesus.HallOrders[floor][Up] = false
		}
		if (!elevExists || cabExists) && (wv.Elevs[wv.ID].CabOrders[floor] == CabO) {
			ordersWithConsesus.CabOrders[floor] = true
		} else {
			ordersWithConsesus.CabOrders[floor] = false
		}
	}
	return ordersWithConsesus
}