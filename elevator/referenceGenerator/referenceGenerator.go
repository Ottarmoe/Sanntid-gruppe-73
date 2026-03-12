package referenceGenerator
import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
)

func ReferenceGenerator(myPhysicalState PhysicalState, myOrders OurOrders) PhysicalState {
	CurrentFloor := myPhysicalState.Floor
	CurrentDirection := myPhysicalState.MovDirection

	switch myPhysicalState.Behaviour {

	case Idle:

		if orderOnCurrentFloorInSameDirection(myPhysicalState, myOrders) {

			referencePhysicalState := setReferencePhysicalState(DoorOpen, CurrentDirection, CurrentFloor)
			return referencePhysicalState
		}

		if orderInSameDirection(myPhysicalState, myOrders) {

			referencePhysicalState := setReferencePhysicalState(Moving, CurrentDirection, CurrentFloor)
			return referencePhysicalState
		}

		if orderOnCurrentFloorInOppositeDirection(myPhysicalState, myOrders) {

			referencePhysicalState := setReferencePhysicalState(DoorOpen, oppositeDirection(CurrentDirection), CurrentFloor)
			return referencePhysicalState
		}

		if orderInOppositeDirection(myPhysicalState, myOrders) {

			referencePhysicalState := setReferencePhysicalState(Moving, oppositeDirection(CurrentDirection), CurrentFloor)
			return referencePhysicalState

		} else {

			return setReferencePhysicalState(Idle, CurrentDirection, CurrentFloor)
		}

	case Moving:

		if  orderOnCurrentFloorInSameDirection(myPhysicalState, myOrders) {
			return setReferencePhysicalState(DoorOpen, CurrentDirection, CurrentFloor)
		}

		if shouldIStopOnNextFloor(myPhysicalState, myOrders) {
			return setReferencePhysicalState(DoorOpen, CurrentDirection, CurrentFloor+directionToIncrement(CurrentDirection))

		} else {
			return setReferencePhysicalState(Moving, CurrentDirection, CurrentFloor+directionToIncrement(CurrentDirection))
		}

	case DoorOpen:
		return setReferencePhysicalState(Idle, CurrentDirection, CurrentFloor)
	}

	return myPhysicalState
}

func setReferencePhysicalState(behavior MotorBehaviour, direction Direction, floor int) PhysicalState {
	var referencePhysicalState PhysicalState
	referencePhysicalState.Behaviour = behavior
	//a reference physical state cannot tell the elevator to move out of the shaft

	if floor < 0 {
		floor = 0
	} else if floor >= NumFloors {
		floor = NumFloors - 1
	}

	referencePhysicalState.Floor = floor

	if floor == 0 {
		referencePhysicalState.MovDirection = Up
	} else if floor == NumFloors-1 {
		referencePhysicalState.MovDirection = Down
	} else {
		referencePhysicalState.MovDirection = direction
	}
	referencePhysicalState.MechError = false

	return referencePhysicalState
}

func directionToIncrement(direction Direction) int {
	if direction == Up {
		return 1
	} else {
		return -1
	}
}

func oppositeDirection(direction Direction) Direction {
	if direction == Up {
		return Down
	} else {
		return Up
	}
}

func orderOnCurrentFloorInSameDirection(me PhysicalState, orders OurOrders) bool {
	floor := me.Floor
	hallOrders := orders.HallOrders
	cabOrders := orders.CabOrders
	direction := me.MovDirection

	return hallOrders[floor][direction] || cabOrders[floor]
}

func orderOnCurrentFloorInOppositeDirection(me PhysicalState, orders OurOrders) bool {
	floor := me.Floor
	hallOrders := orders.HallOrders
	direction := oppositeDirection(me.MovDirection)

	return hallOrders[floor][direction]
}

func orderInSameDirection(me PhysicalState, orders OurOrders) bool {
	return orderInDirection(me.Floor, me.MovDirection, orders)
}

func orderInOppositeDirection(me PhysicalState, orders OurOrders) bool {
	return orderInDirection(me.Floor, oppositeDirection(me.MovDirection), orders)
}

func orderInDirection(currentfloor int, direction Direction, orders OurOrders) bool {
	hallOrders := orders.HallOrders
	cabOrders := orders.CabOrders
	increment := directionToIncrement(direction)

	for floor := currentfloor + increment; floor < NumFloors && floor >= 0; floor += increment {
		if hallOrders[floor][Up] || hallOrders[floor][Down] || cabOrders[floor] {
			return true
		}
	}
	return false
}

func shouldIStopOnNextFloor(me PhysicalState, orders OurOrders) bool {
	direction := me.MovDirection
	nextFloor := me.Floor + directionToIncrement(direction)

	if nextFloor < 0 || nextFloor >= NumFloors {
		return true
	}

	hallOrders := orders.HallOrders
	cabOrders := orders.CabOrders

	if hallOrders[nextFloor][direction] {
		return true
	}
	if cabOrders[nextFloor] {
		return true
	}
	// Stop for opposite-direction order only if there are no orders beyond next floor
	if hallOrders[nextFloor][oppositeDirection(direction)] &&
		!orderInDirection(nextFloor, me.MovDirection, orders) {
		return true
	}
	return false
}
