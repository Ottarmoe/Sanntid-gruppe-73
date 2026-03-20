package referenceGenerator

import (
	. "elevator/elevatorConstants"
	. "elevator/sharedTypes"
)

// Function called on orders assigned by hallRequestAssigner
// to generate a reference next state for the logicalController to steer toward
func ReferenceGenerator(myPhysicalState PhysicalState, myOrders ActionableOrders) PhysicalState {
	CurrentFloor := myPhysicalState.Floor
	CurrentDirection := myPhysicalState.MovDirection

	switch myPhysicalState.Behaviour {

	case Idle:

		if orderOnCurrentFloorInSameDirection(myPhysicalState, myOrders) {
			return setReferencePhysicalState(DoorOpen, CurrentDirection, CurrentFloor)
		}
		if orderInSameDirection(myPhysicalState, myOrders) {
			return setReferencePhysicalState(Moving, CurrentDirection, CurrentFloor)
		}
		if orderOnCurrentFloorInOppositeDirection(myPhysicalState, myOrders) {
			return setReferencePhysicalState(DoorOpen, oppositeDirection(CurrentDirection), CurrentFloor)
		}
		if orderInOppositeDirection(myPhysicalState, myOrders) {
			return setReferencePhysicalState(Moving, oppositeDirection(CurrentDirection), CurrentFloor)
		} else {
			return setReferencePhysicalState(Idle, CurrentDirection, CurrentFloor)
		}

	case Moving:

		if orderOnCurrentFloorInSameDirection(myPhysicalState, myOrders) {
			return setReferencePhysicalState(DoorOpen, CurrentDirection, CurrentFloor)
		}
		if shouldIStopOnNextFloor(myPhysicalState, myOrders) {
			return setReferencePhysicalState(Idle, CurrentDirection, CurrentFloor+directionToIncrement(CurrentDirection))
		} else {
			return setReferencePhysicalState(Moving, CurrentDirection, CurrentFloor+directionToIncrement(CurrentDirection))
		}

	case DoorOpen:
		return setReferencePhysicalState(Idle, CurrentDirection, CurrentFloor)
	}

	return myPhysicalState
}

func setReferencePhysicalState(behavior MotorBehaviour, direction Direction, floor int) PhysicalState {

	//adding safety measures to ensure that the reference physical state does not go out of bounds, which could cause panics in the tests
	if floor <= 0 {
		floor = 0
		direction = Up
	} else if floor >= NumFloors-1 {
		floor = NumFloors - 1
		direction = Down
	}

	var referencePhysicalState PhysicalState
	referencePhysicalState.Behaviour = behavior
	referencePhysicalState.Floor = floor
	referencePhysicalState.MovDirection = direction
	referencePhysicalState.MechError = false

	return referencePhysicalState
}

func directionToIncrement(direction Direction) int {
	if direction == Up {
		return 1
	}
	return -1
}

func oppositeDirection(direction Direction) Direction {
	if direction == Up {
		return Down
	}
	return Up
}

func orderOnCurrentFloorInSameDirection(myPhysicalState PhysicalState, orders ActionableOrders) bool {
	floor := myPhysicalState.Floor
	direction := myPhysicalState.MovDirection
	hallOrders := orders.HallOrders
	cabOrders := orders.CabOrders

	return hallOrders[floor][direction] || cabOrders[floor]
}

func orderOnCurrentFloorInOppositeDirection(myPhysicalState PhysicalState, orders ActionableOrders) bool {
	floor := myPhysicalState.Floor
	direction := oppositeDirection(myPhysicalState.MovDirection)
	hallOrders := orders.HallOrders

	return hallOrders[floor][direction]
}

func orderInSameDirection(myPhysicalState PhysicalState, orders ActionableOrders) bool {
	return orderInDirection(myPhysicalState.Floor, myPhysicalState.MovDirection, orders)
}

func orderInOppositeDirection(myPhysicalState PhysicalState, orders ActionableOrders) bool {
	return orderInDirection(myPhysicalState.Floor, oppositeDirection(myPhysicalState.MovDirection), orders)
}

func orderInDirection(currentfloor int, direction Direction, orders ActionableOrders) bool {
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

func shouldIStopOnNextFloor(myPhysicalState PhysicalState, orders ActionableOrders) bool {
	direction := myPhysicalState.MovDirection
	nextFloor := myPhysicalState.Floor + directionToIncrement(direction)

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
	return hallOrders[nextFloor][oppositeDirection(direction)] &&
		!orderInDirection(nextFloor, myPhysicalState.MovDirection, orders)
}
