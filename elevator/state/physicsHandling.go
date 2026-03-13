package state

import (
	//. "elevator/elevatorConstants"
	. "elevator/stateTypes"
)

func handleFloor(physicalState *PhysicalState, event int) {
	physicalState.Floor = event
}

func handleMotor(wv *ElevWorldView, event PhysicalState) {
	elevator := &wv.ElevStates[wv.ID]
	physics := &elevator.PhysicalState

	physics.Behaviour = event.Behaviour
	physics.MovDirection = event.MovDirection
}

func handleMech(wv *ElevWorldView, event bool) {
	wv.ElevStates[wv.ID].PhysicalState.MechError = event
}

func handleNetworkPhysics(wv *ElevWorldView, netMessage NetMessage) {
	wv.ElevStates[netMessage.ID].PhysicalState = netMessage.ElevState.PhysicalState
}