package state

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
	// "fmt"
)

func handleButton(wv *ElevWorldView, event ButtonEvent) {
	me := wv.MyElev()
	floor := event.Floor

	switch event.Button {
	case BT_HallUp:
		if !wv.AnyoneInHallOrderState(HallOPR, floor, Up) {
			me.OrderState.HallOrders[floor][Up] = HallO
		}
	case BT_HallDown:
		if !wv.AnyoneInHallOrderState(HallOPR, floor, Down) {
			me.OrderState.HallOrders[floor][Down] = HallO
		}
	case BT_Cab:
		if me.OrderState.CabOrders[floor] != CabO {
			me.OrderState.CabOrders[floor] = CabO
			for elev := 0; elev < NumElevators; elev++ {
				wv.CabAgreement[elev][floor] = false
			}
		}
	}
}

func findConsensus(wv ElevWorldView) OrdersWithConsensus {
	var ordersWithConsensus OrdersWithConsensus
	for elev := 0; elev < NumElevators; elev++ {
		for floor := 0; floor < NumFloors; floor++ {
			ordersWithConsensus.CabOrders[elev][floor] = (wv.ElevStates[elev].OrderState.CabOrders[floor] == CabO)
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

			if wv.NetError[elev] == false && ID() != elev {
				anyElevExists = true
				if peerHallOrders[floor][Down] == HallO {
					hallDownExists = true
				}
				if peerHallOrders[floor][Up] == HallO {
					hallUpExists = true
				}
				if wv.CabAgreement[elev][floor] == true {
					cabExists = true
				}
			}
		}
		HallOrders := &wv.ElevStates[ID()].OrderState.HallOrders
		CabOrders := &wv.ElevStates[ID()].OrderState.CabOrders

		if (!anyElevExists || hallDownExists) && (HallOrders[floor][Down] == HallO) {
			ordersWithConsensus.HallOrders[floor][Down] = true
		} else {
			ordersWithConsensus.HallOrders[floor][Down] = false
		}
		if (!anyElevExists || hallUpExists) && (HallOrders[floor][Up] == HallO) {
			ordersWithConsensus.HallOrders[floor][Up] = true
		} else {
			ordersWithConsensus.HallOrders[floor][Up] = false
		}
		if (!anyElevExists || cabExists) && (CabOrders[floor] == CabO) {
			ordersWithConsensus.CabOrders[ID()][floor] = true
		} else {
			ordersWithConsensus.CabOrders[ID()][floor] = false
		}
	}
	return ordersWithConsensus
}

func handleNetworkOrders(wv *ElevWorldView, netMessage NetMessage) {
	wv.ElevStates[netMessage.ID].OrderState.HallOrders = netMessage.ElevState.OrderState.HallOrders

	//atchive their cab state
	for floor := 0; floor < NumFloors; floor++ {
		if netMessage.ElevState.OrderState.CabOrders[floor] != CabUO {
			wv.ElevStates[netMessage.ID].OrderState.CabOrders[floor] = netMessage.ElevState.OrderState.CabOrders[floor]
		}
	}
	//check if their archive is up to date
	for floor := 0; floor < NumFloors; floor++ {
		if netMessage.CabBackups[ID()][floor] == wv.MyElev().OrderState.CabOrders[floor] {
			wv.CabAgreement[netMessage.ID][floor] = true
		} else {
			wv.CabAgreement[netMessage.ID][floor] = false
		}
	}
	//if that elevator not yet seen, integrate their archive
	if !wv.CabArchiveSeen[netMessage.ID] {
		for floor := 0; floor < NumFloors; floor++ {
			if wv.MyElev().OrderState.CabOrders[floor] == CabUO {
				if netMessage.CabBackups[ID()][floor] == CabO {
					wv.MyElev().OrderState.CabOrders[floor] = CabO
				}
			}
		}
		wv.CabArchiveSeen[netMessage.ID] = true
	}
	//if some elevator is in net error, diffuse their cab archive from this other elevator
	for elev := 0; elev < NumElevators; elev++ {
		if wv.NetError[elev] && elev != ID() {
			for floor := 0; floor < NumFloors; floor++ {
				if netMessage.CabBackups[elev][floor] == CabO {
					wv.ElevStates[elev].OrderState.CabOrders[floor] = CabO
				}
			}
		}
	}
}
func handleOrderDynamics(wv *ElevWorldView) {
	myFloor := wv.MyElev().PhysicalState.Floor
	myMovDirection := wv.MyElev().PhysicalState.MovDirection
	myBehaviour := wv.MyElev().PhysicalState.Behaviour

	//transition to OPR when finishing order
	//and removal of cab order
	if myBehaviour == DoorOpen && !wv.MyElev().PhysicalState.MechError {
		//cab order
		//does not need to change CabAgreement
		wv.MyElev().OrderState.CabOrders[myFloor] = CabNO
		//hall order
		if !wv.AnyoneInHallOrderState(HallNO, myFloor, myMovDirection) {
			wv.MyElev().OrderState.HallOrders[myFloor][myMovDirection] = HallOPR
		}
	}
	//order diffusion
	for floor := 0; floor < NumFloors; floor++ {
		for _, dir := range []Direction{Up, Down} {
			orderDiffusion(wv, floor, dir)
		}
	}
}

func orderDiffusion(wv *ElevWorldView, floor int, dir Direction) {
	switch wv.MyElev().OrderState.HallOrders[floor][dir] {
	case HallO:
		if wv.AnyoneInHallOrderState(HallOPR, floor, dir) {
			wv.MyElev().OrderState.HallOrders[floor][dir] = HallOPR
		}
	case HallOPR:
		if !wv.AnyoneInHallOrderState(HallO, floor, dir) {
			wv.MyElev().OrderState.HallOrders[floor][dir] = HallNO
		}
	case HallNO:
		if wv.AnyoneInHallOrderState(HallO, floor, dir) {
			wv.MyElev().OrderState.HallOrders[floor][dir] = HallO
		}
	}
}
