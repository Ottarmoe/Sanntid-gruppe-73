package referenceGenerator

import (
	"elevator/hra"
	. "elevator/state"
	. "elevator/elevatorConstants"
	"fmt"
)

func ReferenceGenerator(
	newState <-chan ElevWorldView,
	//motorRef chan<- PhysicalState,
	//inspRef chan<- PhysicalState,
) {
	for {
		wv := <-newState

		//wv.Elevs[0].NetError = false
		//wv.Elevs[1].NetError = false
		//wv.Elevs[2].NetError = false

		me := wv.Elevs[wv.ID]
		hallOrders := hra.HRA(wv)
		me.HallOrders = hallOrders
		cabOrders := me.CabOrders
		floor := me.CabPhysics.Floor
		behaviour := me.CabPhysics.Motor.Behaviour
		direction := me.CabPhysics.Motor.MovDirection
	

		switch behaviour{
		case Idle:
			anyOrders := orderAfterAndOnFloor(me)



				return 


		case Moving:
			ShouldIStop

		case DoorOpen:
			ShouldICloseDoors
		}

	}
}



func orderSameDirection(me ElevState) bool{
	var increment int;
	direction := me.CabPhysics.Motor.MovDirection

	if direction == Up {
		increment = 1
	}

	if direction == Down {
		increment = -1
	}

	for floor := me.CabPhysics.Floor; floor < NumFloors && floor >= 0; floor+=increment{
		
		if (me.HallOrders[floor][direction] == HallO){
			return true
		}
		if (me.CabOrders[floor] == CabO){
			return true
		}
	}
	return false
}

func orderOnFloor(me ElevState) bool{	
	floor := me.CabPhysics.Floor

	if ((me.HallOrders[floor][Up] == HallO) || (me.HallOrders[floor][Down]) == HallO){
		return true
	}
	if (me.CabOrders[floor] == CabO){
		return true
	}
	return false
}

func orderOppositeDirection(me ElevState) bool{
	var increment int;
	direction := me.CabPhysics.Motor.MovDirection

	if direction == Up {
		increment = -1
	}

	if direction == Down {
		increment = 1
	}

	for floor := me.CabPhysics.Floor; floor < NumFloors && floor >= 0; floor+=increment{
		
		if (me.HallOrders[floor][direction] == HallO){
			return true
		}
		if (me.CabOrders[floor] == CabO){
			return true
		}
	}
	return false
}





