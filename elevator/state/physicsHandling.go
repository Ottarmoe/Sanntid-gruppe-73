package state

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
)

func handleFloor(wv *ElevWorldView, event int) {
	wv.MyElev().PhysicalState.Floor = event
}

func handleMotor(wv *ElevWorldView, event PhysicalState) {
	wv.MyElev().PhysicalState.Behaviour = event.Behaviour
	wv.MyElev().PhysicalState.MovDirection = event.MovDirection
}

func handleMech(wv *ElevWorldView, event bool) {
	wv.MyElev().PhysicalState.MechError = event
}

func handleNetworkPhysics(wv *ElevWorldView, netMessage NetMessage) {
	wv.ElevStates[netMessage.ID].PhysicalState = netMessage.ElevState.PhysicalState
}

func compilePhysicalStates(wv *ElevWorldView) [NumElevators]PhysicalState {
	var completePhysicalState [NumElevators]PhysicalState
	for elev := 0; elev < NumElevators; elev++ {
		completePhysicalState[elev] = wv.ElevStates[elev].PhysicalState
	}
	return completePhysicalState
}
