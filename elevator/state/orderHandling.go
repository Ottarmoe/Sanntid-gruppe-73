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
	ordersWithConsesus.ID = wv.ID
	for elev := 0; elev < NumElevators; elev++ {
		for floor := 0; floor < NumFloors; floor++ {
			ordersWithConsesus.CabOrders[elev][floor] = (wv.ElevStates[elev].OrderState.CabOrders[floor] == CabO)
		}
	}

	for floor := 0; floor < NumFloors; floor++ {
		//check for consensus on _our_ hall and cab order states
		hallDownExists := false
		hallUpExists := false
		cabExists := false
		anyElevExists := false
		for elev := 0; elev < NumElevators; elev++ {
			peerHallOrders := &wv.ElevStates[elev].OrderState.HallOrders

			if wv.NetError[elev] == false && wv.ID != elev {
				anyElevExists = true
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

		if (!anyElevExists || hallDownExists) && (HallOrders[floor][Down] == HallO) {
			ordersWithConsesus.HallOrders[floor][Down] = true
		} else {
			ordersWithConsesus.HallOrders[floor][Down] = false
		}
		if (!anyElevExists || hallUpExists) && (HallOrders[floor][Up] == HallO) {
			ordersWithConsesus.HallOrders[floor][Up] = true
		} else {
			ordersWithConsesus.HallOrders[floor][Up] = false
		}
		if (!anyElevExists || cabExists) && (CabOrders[floor] == CabO) {
			ordersWithConsesus.CabOrders[wv.ID][floor] = true
		} else {
			ordersWithConsesus.CabOrders[wv.ID][floor] = false
		}
	}
	return ordersWithConsesus
}
