package state

import (
	. "elevator/stateTypes"
	. "elevator/elevatorConstants"
)

func handleFloor(physicalState *PhysicalState, event int) {	
	physicalState.Floor = event
}

func handleMotor(wv *ElevWorldView, event PhysicalState) {
	elevator := &wv.ElevStates[wv.ID]	
	physics := &elevator.PhysicalState

	physics.Behaviour = event.Behaviour
	if event.Behaviour == Moving {
		physics.MovDirection = event.MovDirection
	}
}

func handleMech(wv *ElevWorldView, event bool) {
	wv.ElevStates[wv.ID].PhysicalState.MechError = event
}

func handleNetworkPhysics(wv *ElevWorldView, netMessage NetMessage) {
	for elev := 0; elev < NumElevators; elev++ {
		if elev != wv.ID {
			wv.ElevStates[elev].PhysicalState = netMessage.ElevState.PhysicalState
		}
	}
}