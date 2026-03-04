package state

func handleFloor(wv *ElevWorldView, event int) {
	wv.Elevs[wv.ID].CabPhysics.Floor = event
}

func handleMotor(wv *ElevWorldView, event MotorState) {
	wv.Elevs[wv.ID].CabPhysics.Motor.Behaviour = event.Behaviour
	if event.Behaviour == Moving {
		wv.Elevs[wv.ID].CabPhysics.Motor.MovDirection = event.MovDirection
	}
}

func handleMech(wv *ElevWorldView, event bool) {
	wv.Elevs[wv.ID].CabMechError = event
}