package state

func handleFloor(physicalState *PhysicalState, event int) {	
	physicalState.Floor = event
}

func handleMotor(wv *ElevWorldView, event PhysicalState) {
	elevator := &wv.Elevs[wv.ID]	
	physics := &elevator.PhysicalState

	physics.Behaviour = event.Behaviour
	if event.Behaviour == Moving {
		physics.MovDirection = event.MovDirection
	}
}

func handleMech(wv *ElevWorldView, event bool) {
	wv.Elevs[wv.ID].MechError = event
}