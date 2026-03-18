package state

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
	. "elevio"
	// "fmt"
)

func handleButton(wv *ElevWorldView, event ButtonEvent) {
	me := &wv.ElevStates[ID()]
	switch event.Button {
	case BT_HallUp:
		//is anyone in HallOPR?
		readyToTransition := true
		for elev := 0; elev < NumElevators; elev++ {
			if !wv.NetError[elev] || elev == ID() {
				if wv.ElevStates[elev].OrderState.HallOrders[event.Floor][Up] == HallOPR {
					readyToTransition = false
				}
			}
		}
		if readyToTransition {
			me.OrderState.HallOrders[event.Floor][Up] = HallO
		}
	case BT_HallDown:
		//is anyone in HallOPR?
		readyToTransition := true
		for elev := 0; elev < NumElevators; elev++ {
			if !wv.NetError[elev] || elev == ID() {
				if wv.ElevStates[elev].OrderState.HallOrders[event.Floor][Down] == HallOPR {
					readyToTransition = false
				}
			}
		}
		if readyToTransition {
			me.OrderState.HallOrders[event.Floor][Down] = HallO
		}
	case BT_Cab:
		if me.OrderState.CabOrders[event.Floor] != CabO {
			me.OrderState.CabOrders[event.Floor] = CabO
			for elev := 0; elev < NumElevators; elev++ {
				wv.CabAgreement[elev][event.Floor] = false
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
		if netMessage.CabBackups[ID()][floor] == wv.ElevStates[ID()].OrderState.CabOrders[floor] {
			wv.CabAgreement[netMessage.ID][floor] = true
		} else {
			wv.CabAgreement[netMessage.ID][floor] = false
		}
	}
	//if that elevator not yet seen, integrate their archive
	if !wv.CabArchiveSeen[netMessage.ID] {
		for floor := 0; floor < NumFloors; floor++ {
			if wv.ElevStates[ID()].OrderState.CabOrders[floor] == CabUO {
				if netMessage.CabBackups[ID()][floor] == CabO {
					wv.ElevStates[ID()].OrderState.CabOrders[floor] = CabO
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
	elevator := &wv.ElevStates[ID()]
	physics := &elevator.PhysicalState

	//transition to OPR when finishing order
	//and removal of cab order
	if physics.Behaviour == DoorOpen && !elevator.PhysicalState.MechError {
		//cab order
		//does not need to change CabAgreement
		elevator.OrderState.CabOrders[physics.Floor] = CabNO
		//hall order
		//is everyone in Order or OPR?
		readyToTransition := true
		for elev := 0; elev < NumElevators; elev++ {
			if !wv.NetError[elev] || elev == ID() {
				readyToTransition = readyToTransition && wv.ElevStates[elev].OrderState.HallOrders[physics.Floor][physics.MovDirection] != HallNO
			}
		}
		if readyToTransition {
			elevator.OrderState.HallOrders[physics.Floor][physics.MovDirection] = HallOPR
		}
	}
	//order diffusion
	for floor := 0; floor < NumFloors; floor++ {
		for _, dir := range []Direction{Up, Down} {
			//construct array of other elevators
			elevatorStates := []HallOrderState{}
			for elev := 0; elev < NumElevators; elev++ {
				if !wv.NetError[elev] || elev == ID() {
					elevatorStates = append(elevatorStates, wv.ElevStates[elev].OrderState.HallOrders[floor][dir])
				}
			}
			elevator.OrderState.HallOrders[floor][dir] = SingleOrderDiffusion(elevator.OrderState.HallOrders[floor][dir], elevatorStates)
		}
	}
}

func SingleOrderDiffusion(me HallOrderState, orderStates []HallOrderState) HallOrderState {
	switch me {
	case HallO:
		//transition to OPR if someone is in OPR
		readyToTransition := false
		for _, stat := range orderStates {
			if stat == HallOPR {
				readyToTransition = true
			}
		}
		if readyToTransition {
			return HallOPR
		}
	case HallOPR:
		//transition to HallNO if no one is in HallO anymore
		readyToTransition := true
		for _, stat := range orderStates {
			if stat == HallO {
				readyToTransition = false
			}
		}
		if readyToTransition {
			return HallNO
		}
	case HallNO:
		//transition to HallO if someone is in HallO
		readyToTransition := false
		for _, stat := range orderStates {
			if stat == HallO {
				readyToTransition = true
			}
		}
		if readyToTransition {
			return HallO
		}
	}
	return me
}
