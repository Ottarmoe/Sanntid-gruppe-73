package referenceGenerator

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
)

func ReferenceGenerator(myPhysicalState PhysicalState, myOrders AssignedOrders) PhysicalState {
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

func orderOnCurrentFloorInSameDirection(me PhysicalState, orders AssignedOrders) bool {
	floor := me.Floor
	direction := me.MovDirection
	hallOrders := orders.HallOrders
	cabOrders := orders.CabOrders

	return hallOrders[floor][direction] || cabOrders[floor]
}

func orderOnCurrentFloorInOppositeDirection(me PhysicalState, orders AssignedOrders) bool {
	floor := me.Floor
	direction := oppositeDirection(me.MovDirection)
	hallOrders := orders.HallOrders

	return hallOrders[floor][direction]
}

func orderInSameDirection(me PhysicalState, orders AssignedOrders) bool {
	return orderInDirection(me.Floor, me.MovDirection, orders)
}

func orderInOppositeDirection(me PhysicalState, orders AssignedOrders) bool {
	return orderInDirection(me.Floor, oppositeDirection(me.MovDirection), orders)
}

func orderInDirection(currentfloor int, direction Direction, orders AssignedOrders) bool {
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

func shouldIStopOnNextFloor(me PhysicalState, orders AssignedOrders) bool {
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
	return hallOrders[nextFloor][oppositeDirection(direction)] &&
		!orderInDirection(nextFloor, me.MovDirection, orders)
}
