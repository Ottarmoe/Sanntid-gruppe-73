package hardware

import (
	. "elevator/elevatorConstants"
	. "elevator/hardwareLow"
	. "elevator/sharedTypes"
	"time"
)

const pollRate = 20 * time.Millisecond

// Listens for changes in physical state or button lights communicated by statekeeper,
// and controls hardware to make sure that the physical world reflects the worldview.
func HardwareOut(physicalToHardwareCh <-chan PhysicalState, ordersWithConsensusToHardwareCh <-chan OrdersWithConsensus) {
	var prevConsensus OrdersWithConsensus
	resetLights()
	for {
		select {
		case physicalState := <-physicalToHardwareCh:
			updateMotorAndDoor(physicalState)

		case ordersWithConsensus := <-ordersWithConsensusToHardwareCh:
			updateButtonLights(ordersWithConsensus, prevConsensus)
			prevConsensus = ordersWithConsensus
		}
	}
}

// *HardwareIn*
// Three poll routines that continuously read hardware inputs (buttons, floor sensor,
// and obstruction switch) and send events to statekeeper only
// when a change is detected.
func PollButtons(receiverCh chan<- ButtonEvent) {
	prev := make([][3]bool, NumFloors)
	for {
		time.Sleep(pollRate)
		for f := 0; f < NumFloors; f++ {
			for b := ButtonType(0); b < 3; b++ {
				v := GetButton(b, f)
				if v != prev[f][b] && v != false {
					receiverCh <- ButtonEvent{Floor: f, Button: ButtonType(b)}
				}
				prev[f][b] = v
			}
		}
	}
}
func PollFloorSensor(receiverCh chan<- int) {
	prev := -1
	for {
		time.Sleep(pollRate)
		v := GetFloor()
		if v != prev && v != -1 {
			receiverCh <- v
		}
		prev = v
	}
}
func PollObstructionSwitch(receiverCh chan<- bool) {
	prev := false
	for {
		time.Sleep(pollRate)
		v := GetObstruction()
		if v != prev {
			receiverCh <- v
		}
		prev = v
	}
}

// updateMotorAndDoor drives motor and door lamp to match the current behaviour
func updateMotorAndDoor(physicalState PhysicalState) {
	SetFloorIndicator(physicalState.Floor)
	switch physicalState.Behaviour {
	case Idle:
		SetMotorDirection(MotorDirStop)
		SetDoorOpenLamp(false)
	case Moving:
		SetDoorOpenLamp(false)
		if physicalState.MovDirection == Up {
			SetMotorDirection(MotorDirUp)
		}
		if physicalState.MovDirection == Down {
			SetMotorDirection(MotorDirDown)
		}
	case DoorOpen:
		SetMotorDirection(MotorDirStop)
		SetDoorOpenLamp(true)
	}
	//bonus feature :D
	//mark mech error by lighting the stop button
	SetStopLamp(physicalState.MechError)
}

// updateButtonLights updates button lamps to reflect agreed-upon orders.
// Only updates lamps that have changed since the last update to limit redundant hardware call
func updateButtonLights(orders OrdersWithConsensus, prev OrdersWithConsensus) {
	for floor := 0; floor < NumFloors; floor++ {
		if orders.HallOrders[floor][Down] != prev.HallOrders[floor][Down] {
			SetButtonLamp(ButtonHallDown, floor, orders.HallOrders[floor][Down])
		}
		if orders.HallOrders[floor][Up] != prev.HallOrders[floor][Up] {
			SetButtonLamp(ButtonHallUp, floor, orders.HallOrders[floor][Up])
		}
		if orders.CabOrders[MyID()][floor] != prev.CabOrders[MyID()][floor] {
			SetButtonLamp(ButtonCab, floor, orders.CabOrders[MyID()][floor])
		}
	}
}

func resetLights() {
	for floor := 0; floor < NumFloors; floor++ {
		SetButtonLamp(ButtonHallDown, floor, false)
		SetButtonLamp(ButtonHallUp, floor, false)
		SetButtonLamp(ButtonCab, floor, false)
	}
}
