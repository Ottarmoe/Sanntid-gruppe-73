package referenceGenerator

import (
	//"elevator/hra"
	. "elevator/state"
	. "elevator/elevatorConstants"
)

func ReferenceGenerator(newState <-chan ElevWorldView) PhysicalState {
	for {
		newState := <- newState
		myElevState := newState.ElevStates[newState.ID]

		CurrentFloor := myElevState.PhysicalState.Floor
		CurrentDirection := myElevState.PhysicalState.MovDirection

	
		switch myElevState.PhysicalState.Behaviour{

		case Idle:

			anyOrdersOnFloorInSameDirection := orderOnCurrentFloorInSameDirection(myElevState)
			if anyOrdersOnFloorInSameDirection {
				
				referencePhysicalState := setReferencePhysicalState(DoorOpen, CurrentDirection, CurrentFloor)
				return referencePhysicalState
			}
			

			anyOrdersInSameDirection := orderInSameDirection(myElevState)
			if anyOrdersInSameDirection {

				referencePhysicalState := setReferencePhysicalState(Moving, CurrentDirection, CurrentFloor)
				return referencePhysicalState
			}

			anyOrdersOnFloorInOppositeDirection := orderOnCurrentFloorInOppositeDirection(myElevState)
			if anyOrdersOnFloorInOppositeDirection {

				referencePhysicalState := setReferencePhysicalState(DoorOpen, oppositeDirection(CurrentDirection), CurrentFloor)
				return referencePhysicalState
			}

			anyOrdersInOppositeDirection := orderInOppositeDirection(myElevState)
			if anyOrdersInOppositeDirection {


				referencePhysicalState := setReferencePhysicalState(Moving, oppositeDirection(CurrentDirection), CurrentFloor)
				return referencePhysicalState

			} else{

				referencePhysicalState := setReferencePhysicalState(Idle, CurrentDirection, CurrentFloor)
				return referencePhysicalState
				} 
		
		case Moving:

			shouldIStop := ShouldIStopOnNextFloor(myElevState)
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

	}
}


func setReferencePhysicalState(behavior MotorBehaviour, direction Direction, floor int) PhysicalState{
	var referencePhysicalState PhysicalState
	referencePhysicalState.Behaviour = behavior
	referencePhysicalState.MovDirection = direction
	referencePhysicalState.Floor = floor
	referencePhysicalState.MechError = false

	return referencePhysicalState
}


func directionToIncrement(direction Direction) int{
	if direction == Up {
		return 1
	}
	if direction == Down {
		return -1
	}
	return 0
}

func oppositeDirection(direction Direction) Direction{
	if direction == Up {
		return Down
	}
	if direction == Down {
		return Up
	}
	return 0
}


func orderOnCurrentFloorInSameDirection(me ElevState) bool{
	floor := me.PhysicalState.Floor
	hallOrders := me.OrderState.HallOrders
	cabOrders := me.OrderState.CabOrders
	direction := me.PhysicalState.MovDirection
	
	if ((hallOrders[floor][direction] == HallO) || cabOrders[floor] == CabO){
		return true
	}
	return false
}

func orderOnCurrentFloorInOppositeDirection(me ElevState) bool{
	floor := me.PhysicalState.Floor
	hallOrders := me.OrderState.HallOrders
	direction := me.PhysicalState.MovDirection *-1
	
	if (hallOrders[floor][direction] == HallO){
		return true
	}
	return false
}

func orderInSameDirection(me ElevState) bool{
	direction := me.PhysicalState.MovDirection
	currentfloor := me.PhysicalState.Floor
	hallOrders := me.OrderState.HallOrders
	cabOrders := me.OrderState.CabOrders
	increment := directionToIncrement(direction)

	for floor := currentfloor; floor < NumFloors && floor >= 0; floor+=increment{
		
		if (hallOrders[floor][direction] == HallO){
			return true
		}
		if (cabOrders[floor] == CabO){
			return true
		}
	}
	return false
}


func orderInOppositeDirection(me ElevState) bool{
	direction := me.PhysicalState.MovDirection
	currentfloor := me.PhysicalState.Floor
	hallOrders := me.OrderState.HallOrders
	cabOrders := me.OrderState.CabOrders
	increment := directionToIncrement(direction) * -1


	for floor := currentfloor; floor < NumFloors && floor >= 0; floor+=increment{
		
		if (hallOrders[floor][direction] == HallO){
			return true
		}
		if (cabOrders[floor] == CabO){
			return true
		}
	}
	return false
}


func ShouldIStopOnNextFloor(me ElevState) bool{
	direction := me.PhysicalState.MovDirection
	nextFloor := me.PhysicalState.Floor + directionToIncrement(direction)
	hallOrders := me.OrderState.HallOrders
	cabOrders := me.OrderState.CabOrders

	if (hallOrders[nextFloor][direction] == HallO){
		return true
	}
	if (cabOrders[nextFloor] == CabO){
		return true
	}
	return false
}
