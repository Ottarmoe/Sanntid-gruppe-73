package referenceGenerator

import (
	. "elevator/hra"
	. "elevator/state"
	. "elevator/elevatorConstants"
)

func ReferenceGenerator(physicalState <-chan PhysicalState, ourOrders <-chan OurOrders) PhysicalState {
//TO DO: change hallorders to the assigned hall orders, not the ones from the worldview. 
// The hall orders from the worldview are the ones that have not been assigned to an elevator yet,
//  but the reference generator should use the assigned hall orders, which are in the orderstate of the elevstate.
	for {
		myPhysicalState := <-physicalState
		myOrders := <-ourOrders
		
		CurrentFloor := myPhysicalState.Floor
		CurrentDirection := myPhysicalState.MovDirection

	
		switch myPhysicalState.Behaviour{

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

			} else{

				referencePhysicalState := setReferencePhysicalState(Idle, CurrentDirection, CurrentFloor)
				return referencePhysicalState
				} 
		
		case Moving:

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


func orderOnCurrentFloorInSameDirection(me PhysicalState, orders OurOrders) bool{
	floor := me.Floor
	hallOrders := orders.HallOrders
	cabOrders := orders.CabOrders
	direction := me.MovDirection
	
	if ((hallOrders[floor][direction]) || cabOrders[floor]){
		return true
	}
	return false
}

func orderOnCurrentFloorInOppositeDirection(me PhysicalState, orders OurOrders) bool{
	floor := me.Floor
	hallOrders := orders.HallOrders
	direction := oppositeDirection(me.MovDirection)
	
	if (hallOrders[floor][direction]){
		return true
	}
	return false
}

func orderInSameDirection(me PhysicalState, orders OurOrders) bool{
	direction := me.MovDirection
	currentfloor := me.Floor
	hallOrders := orders.HallOrders
	cabOrders := orders.CabOrders
	increment := directionToIncrement(direction)

	for floor := currentfloor; floor < NumFloors && floor >= 0; floor+=increment{
		
		if (hallOrders[floor][Up] || hallOrders[floor][Down] || cabOrders[floor]){
			return true
		}
	}
	return false
}


func orderInOppositeDirection(me PhysicalState, orders OurOrders) bool{
	direction := me.MovDirection
	currentfloor := me.Floor
	hallOrders := orders.HallOrders
	cabOrders := orders.CabOrders
	increment := directionToIncrement(direction) * -1


	for floor := currentfloor; floor < NumFloors && floor >= 0; floor+=increment{
		
		if (hallOrders[floor][Up] || hallOrders[floor][Down] || cabOrders[floor]){
			return true
		}
	}
	return false
}


func ShouldIStopOnNextFloor(me PhysicalState, orders OurOrders) bool{
	direction := me.MovDirection
	nextFloor := me.Floor + directionToIncrement(direction)
	hallOrders := orders.HallOrders
	cabOrders := orders.CabOrders

	if (hallOrders[nextFloor][direction]) {
		return true
	}
	if (cabOrders[nextFloor]){
		return true
	}
	if (hallOrders[nextFloor][oppositeDirection(direction)] && !orderInSameDirection(me, orders)){
		return true
	}
		return false
}
