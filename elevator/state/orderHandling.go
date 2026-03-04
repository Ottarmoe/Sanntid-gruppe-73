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

func findConsensus(wv *ElevWorldView) ElevWorldView {
	consensus := *wv
	for floor := 0; floor < NumFloors; floor++ {
		hallDownExists := false
		hallUpExists := false
		cabExists := false
		elevExists := false
		for elev := 0; elev < NumElevators; elev++ {
			if wv.Elevs[elev].NetError == false {
				elevExists = true
				if wv.Elevs[elev].HallOrders[floor][Up] == HallO {
					hallDownExists = true
				}
				if wv.Elevs[elev].HallOrders[floor][Down] == HallO {
					hallUpExists = true
				}
				if wv.Elevs[elev].CabAgreement[floor] == true {
					cabExists = true
				}
			}
		}
		if (!elevExists || hallDownExists) && (wv.Elevs[wv.ID].HallOrders[floor][Down] == HallO) {
			consensus.Elevs[consensus.ID].HallOrders[floor][Down] = HallO
		} else {
			consensus.Elevs[consensus.ID].HallOrders[floor][Down] = HallNO
		}
		if (!elevExists || hallUpExists) && (wv.Elevs[wv.ID].HallOrders[floor][Up] == HallO) {
			consensus.Elevs[consensus.ID].HallOrders[floor][Up] = HallO
		} else {
			consensus.Elevs[consensus.ID].HallOrders[floor][Up] = HallNO
		}
		if (!elevExists || cabExists) && (wv.Elevs[wv.ID].CabOrders[floor] == CabO) {
			consensus.Elevs[consensus.ID].CabOrders[floor] = CabO
		} else {
			consensus.Elevs[consensus.ID].CabOrders[floor] = CabNO
		}
	}
	return consensus
}