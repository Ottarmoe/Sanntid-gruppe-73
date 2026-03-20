package hardware

import (
	. "elevator/elevatorConstants"
	. "elevator/hardwareLow"
	. "elevator/stateTypes"
	"time"
)

const pollRate = 20 * time.Millisecond

// Listens for changes in physical state or button lights communicated by statekeeper,
// and cotrols hardware to make sure that the physical world reflects the worldview.
func HardwareOut(physicalToHardware <-chan PhysicalState, ordersWithConsensusToHardware <-chan OrdersWithConsensus) {
	var prevConsensus OrdersWithConsensus
	resetLights()
	for {
		select {
		case physicalState := <-physicalToHardware:
			updateMotorAndDoors(physicalState)

		case ordersWithConsensus := <-ordersWithConsensusToHardware:
			updateButtonLights(ordersWithConsensus, prevConsensus)
			prevConsensus = ordersWithConsensus
		}
	}
}

// *HardwareIn*
// Four poll routines that continuously read hardware inputs (buttons, floor sensor,
// stop button, obstruction switch) and send events to statekeeper only
// when a change is detected.
func PollButtons(receiver chan<- ButtonEvent) {
	prev := make([][3]bool, NumFloors)
	for {
		time.Sleep(pollRate)
		for f := 0; f < NumFloors; f++ {
			for b := ButtonType(0); b < 3; b++ {
				v := GetButton(b, f)
				if v != prev[f][b] && v != false {
					receiver <- ButtonEvent{Floor: f, Button: ButtonType(b)}
				}
				prev[f][b] = v
			}
		}
	}
}
func PollFloorSensor(receiver chan<- int) {
	prev := -1
	for {
		time.Sleep(pollRate)
		v := GetFloor()
		if v != prev && v != -1 {
			receiver <- v
		}
		prev = v
	}
}
func PollStopButton(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(pollRate)
		v := GetStop()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}
func PollObstructionSwitch(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(pollRate)
		v := GetObstruction()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

// updateMotorAndDoors drives motor and door lamp to match the current behaviour
func updateMotorAndDoors(physicalState PhysicalState) {
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
}

// updateButtonLights updates button lamps to reflect agreed-upon orders.
// Only updates lamps that have changed since the last update to minimize redundant hardware call
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
