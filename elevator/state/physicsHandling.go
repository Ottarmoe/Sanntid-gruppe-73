package state

import (
	. "elevator/stateTypes"
)

func handleFloor(wv *ElevWorldView, event int) {
	me := wv.MyState()

	me.PhysicalState.Floor = event
}

func handleMotor(wv *ElevWorldView, event PhysicalState) {
	me := wv.MyState()

	me.PhysicalState.Behaviour = event.Behaviour
	me.PhysicalState.MovDirection = event.MovDirection
}

func handleMech(wv *ElevWorldView, event bool) {
	me := wv.MyState()

	me.PhysicalState.MechError = event
}

func handleNetworkPhysics(wv *ElevWorldView, netMessage NetMessage) {
	wv.ElevStates[netMessage.ID].PhysicalState = netMessage.ElevState.PhysicalState
}