package referenceGenerator

import (
	//. "elevator/hra"
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
)

func ReferenceGenerator(myPhysicalState PhysicalState, myOrders OurOrders) PhysicalState {
	//TO DO: change hallorders to the assigned hall orders, not the ones from the worldview.
	// The hall orders from the worldview are the ones that have not been assigned to an elevator yet,
	//  but the reference generator should use the assigned hall orders, which are in the orderstate of the elevstate.
	//...is this still relevant?
	CurrentFloor := myPhysicalState.Floor
	CurrentDirection := myPhysicalState.MovDirection

	switch myPhysicalState.Behaviour {

	case Idle:

		anyOrdersOnFloorInSameDirection := orderOnCurrentFloorInSameDirection(myPhysicalState, myOrders)
		if anyOrdersOnFloorInSameDirection {

			referencePhysicalState := setReferencePhysicalState(DoorOpen, CurrentDirection, CurrentFloor)
			return referencePhysicalState
		}

		anyOrdersInSameDirection := orderInSameDirection(myPhysicalState, myOrders)
		if anyOrdersInSameDirection {

			referencePhysicalState := setReferencePhysicalState(Moving, CurrentDirection, CurrentFloor)
			return referencePhysicalState
		}

		anyOrdersOnFloorInOppositeDirection := orderOnCurrentFloorInOppositeDirection(myPhysicalState, myOrders)
		if anyOrdersOnFloorInOppositeDirection {

			referencePhysicalState := setReferencePhysicalState(DoorOpen, oppositeDirection(CurrentDirection), CurrentFloor)
			return referencePhysicalState
		}

		anyOrdersInOppositeDirection := orderInOppositeDirection(myPhysicalState, myOrders)
		if anyOrdersInOppositeDirection {

			referencePhysicalState := setReferencePhysicalState(Moving, oppositeDirection(CurrentDirection), CurrentFloor)
			return referencePhysicalState

		} else {

			referencePhysicalState := setReferencePhysicalState(Idle, CurrentDirection, CurrentFloor)
			return referencePhysicalState
		}

	case Moving:

		anyOrdersOnFloorInSameDirection := orderOnCurrentFloorInSameDirection(myPhysicalState, myOrders)
		if anyOrdersOnFloorInSameDirection {

			referencePhysicalState := setReferencePhysicalState(DoorOpen, CurrentDirection, CurrentFloor)
			return referencePhysicalState
		}

		shouldIStop := ShouldIStopOnNextFloor(myPhysicalState, myOrders)
		if shouldIStop {
			referencePhysicalState := setReferencePhysicalState(DoorOpen, CurrentDirection, CurrentFloor+directionToIncrement(CurrentDirection))
			return referencePhysicalState

		} else {
			referencePhysicalState := setReferencePhysicalState(Moving, CurrentDirection, CurrentFloor+directionToIncrement(CurrentDirection))
			return referencePhysicalState
		}

	case DoorOpen:
		referencePhysicalState := setReferencePhysicalState(Idle, CurrentDirection, CurrentFloor)
		return referencePhysicalState
	}

	return myPhysicalState
}

func setReferencePhysicalState(behavior MotorBehaviour, direction Direction, floor int) PhysicalState {
	var referencePhysicalState PhysicalState
	referencePhysicalState.Behaviour = behavior
	//a reference physical state cannot tell the elevator to move out of the shaft
	if floor == 0 {
		referencePhysicalState.MovDirection = Up
	} else if floor == NumFloors-1 {
		referencePhysicalState.MovDirection = Down
	} else {
		referencePhysicalState.MovDirection = direction
	}
	referencePhysicalState.Floor = floor
	referencePhysicalState.MechError = false

	return referencePhysicalState
}

func directionToIncrement(direction Direction) int {
	if direction == Up {
		return 1
	}
	if direction == Down {
		return -1
	}
	return 0
}

func oppositeDirection(direction Direction) Direction {
	if direction == Up {
		return Down
	}
	if direction == Down {
		return Up
	}
	return 0
}

func orderOnCurrentFloorInSameDirection(me PhysicalState, orders OurOrders) bool {
	floor := me.Floor
	hallOrders := orders.HallOrders
	cabOrders := orders.CabOrders
	direction := me.MovDirection

	if (hallOrders[floor][direction]) || cabOrders[floor] {
		return true
	}
	return false
}

func orderOnCurrentFloorInOppositeDirection(me PhysicalState, orders OurOrders) bool {
	floor := me.Floor
	hallOrders := orders.HallOrders
	direction := oppositeDirection(me.MovDirection)

	if hallOrders[floor][direction] {
		return true
	}
	return false
}

func orderInSameDirection(me PhysicalState, orders OurOrders) bool {
	return orderInDirection(me.Floor, me.MovDirection, orders)
}

func orderInOppositeDirection(me PhysicalState, orders OurOrders) bool {
	return orderInDirection(me.Floor, oppositeDirection(me.MovDirection), orders)
}

func orderInDirection(currentfloor int, direction Direction, orders OurOrders) bool {
	//return directionalOrderInDirection(me.MovDirection, me.MovDirection, me.Floor, orders)
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

func ShouldIStopOnNextFloor(me PhysicalState, orders OurOrders) bool {
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
