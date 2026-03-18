package hardwareControl

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
	. "elevio"
	"time"
)

const pollRate = 20 * time.Millisecond

func HardWareControl(physicalToHardware <-chan PhysicalState, ordersWithConsensusToHardwarech <-chan OrdersWithConsensus) {
	var prevConsensus OrdersWithConsensus
	resetLights()
	for {
		select {
		case physicalState := <-physicalToHardware:
			physicalStateToHardware(physicalState)

		case ordersWithConsensus := <-ordersWithConsensusToHardwarech:
			ordersWithConsensusToHardware(ordersWithConsensus, prevConsensus)
			prevConsensus = ordersWithConsensus
		}
	}
}

// physicalStateToHardware drives motor and door lamp to match the current behaviour
func physicalStateToHardware(state PhysicalState) {
	SetFloorIndicator(state.Floor)
	switch state.Behaviour {
	case Idle:
		SetMotorDirection(MD_Stop)
		SetDoorOpenLamp(false)
	case Moving:
		SetDoorOpenLamp(false)
		if state.MovDirection == Up {
			SetMotorDirection(MD_Up)
		}
		if state.MovDirection == Down {
			SetMotorDirection(MD_Down)
		}
	case DoorOpen:
		SetMotorDirection(MD_Stop)
		SetDoorOpenLamp(true)
	}
}

// ordersWithConsensusToHardware updates button lamps to reflect agreed-upon orders.
// Only updates lamps that have changed since the last update to minimize redundant hardware call
func ordersWithConsensusToHardware(orders OrdersWithConsensus, prev OrdersWithConsensus) {
	for floor := 0; floor < NumFloors; floor++ {
		if orders.HallOrders[floor][Down] != prev.HallOrders[floor][Down] {
			SetButtonLamp(BT_HallDown, floor, orders.HallOrders[floor][Down])
		}
		if orders.HallOrders[floor][Up] != prev.HallOrders[floor][Up] {
			SetButtonLamp(BT_HallUp, floor, orders.HallOrders[floor][Up])
		}
		if orders.CabOrders[ID()][floor] != prev.CabOrders[ID()][floor] {
			SetButtonLamp(BT_Cab, floor, orders.CabOrders[ID()][floor])
		}
	}
}

func resetLights() {
	for floor := 0; floor < NumFloors; floor++ {
		SetButtonLamp(BT_HallDown, floor, false)
		SetButtonLamp(BT_HallUp, floor, false)
		SetButtonLamp(BT_Cab, floor, false)
	}
}

func PollButtons(receiver chan<- ButtonEvent) {
	prev := make([][3]bool, NumFloors)
	for {
		time.Sleep(pollRate)
		for f := 0; f < NumFloors; f++ {
			for b := ButtonType(0); b < 3; b++ {
				v := GetButton(b, f)
				if v != prev[f][b] && v != false {
					receiver <- ButtonEvent{f, ButtonType(b)}
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