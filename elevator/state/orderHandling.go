package state

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
)

func handleButton(wv *ElevWorldView, event ButtonEvent) {
	me := wv.MyElev()
	floor := event.Floor

	switch event.Button {
	case ButtonHallUp:
		if !wv.AnyoneInHallOrderState(HallOPR, floor, Up) {
			me.OrderState.HallOrders[floor][Up] = HallO
		}
	case ButtonHallDown:
		if !wv.AnyoneInHallOrderState(HallOPR, floor, Down) {
			me.OrderState.HallOrders[floor][Down] = HallO
		}
	case ButtonCab:
		if me.OrderState.CabOrders[floor] != CabO {
			me.OrderState.CabOrders[floor] = CabO
			for elev := 0; elev < NumElevators; elev++ {
				wv.CabAgreement[elev][floor] = false
			}
		}
	}
}

func findConsensus(wv *ElevWorldView) OrdersWithConsensus {
	var ordersWithConsensus OrdersWithConsensus
	//copy the cab orders from all elevators
	for elev := 0; elev < NumElevators; elev++ {
		for floor := 0; floor < NumFloors; floor++ {
			ordersWithConsensus.CabOrders[elev][floor] = (wv.ElevStates[elev].OrderState.CabOrders[floor] == CabO)
		}
	}

	//check for consensus on our hall and cab order states
	for floor := 0; floor < NumFloors; floor++ {
		//hall orders
		for _, dir := range []Direction{Up, Down} {
			if !wv.AnyPeerExists() {
				if wv.MyElev().OrderState.HallOrders[floor][dir] == HallO {
					ordersWithConsensus.HallOrders[floor][dir] = true
				}
			} else if !wv.AnyoneInHallOrderState(HallOPR, floor, dir) &&
				wv.AnyoneElseInHallOrderState(HallO, floor, dir) {
				ordersWithConsensus.HallOrders[floor][dir] = true
			} else {
				ordersWithConsensus.HallOrders[floor][dir] = false
			}
		}
		//negate cab order if not archived
		if wv.AnyPeerExists() && !wv.CabOrderArchiveExists(floor) {
			ordersWithConsensus.CabOrders[MyID()][floor] = false
		}
	}
	return ordersWithConsensus
}

func handleNetworkOrders(wv *ElevWorldView, netMessage NetMessage) {
	wv.ElevStates[netMessage.ID].OrderState.HallOrders = netMessage.ElevState.OrderState.HallOrders

	//archive their cab state
	for floor := 0; floor < NumFloors; floor++ {
		if netMessage.ElevState.OrderState.CabOrders[floor] != CabUO {
			wv.ElevStates[netMessage.ID].OrderState.CabOrders[floor] = netMessage.ElevState.OrderState.CabOrders[floor]
		}
	}
	//if that elevator not yet seen, integrate their archive
	if !wv.CabArchiveSeen[netMessage.ID] {
		//fmt.Println("integrating new archive")
		for floor := 0; floor < NumFloors; floor++ {
			if netMessage.CabBackups[MyID()][floor] == CabO {
				wv.MyElev().OrderState.CabOrders[floor] = CabO
			}
		}
		wv.CabArchiveSeen[netMessage.ID] = true
	}
	//check if their archive is up to date
	for floor := 0; floor < NumFloors; floor++ {
		if netMessage.CabBackups[MyID()][floor] == wv.MyElev().OrderState.CabOrders[floor] {
			wv.CabAgreement[netMessage.ID][floor] = true
		} else {
			wv.CabAgreement[netMessage.ID][floor] = false
		}
	}
	//if some other elevator is not online, accept a copy of the archive of this elevator cab order state from the message
	for elev := 0; elev < NumElevators; elev++ {
		if wv.IsOfflinePeer(elev) {
			for floor := 0; floor < NumFloors; floor++ {
				if netMessage.CabBackups[elev][floor] == CabO {
					wv.ElevStates[elev].OrderState.CabOrders[floor] = CabO
				}
			}
		}
	}
}

// close orders if the elevator is in the right state to service them
// Update your hall orders if your image of another elevator has a more recent order state
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
