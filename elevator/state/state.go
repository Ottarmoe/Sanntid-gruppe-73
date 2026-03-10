package state

import (
	. "elevator/elevatorConstants"
	"elevator/hra"
	"elevator/referenceGenerator"
	. "elevator/stateTypes"
	. "elevio"
	"fmt"
)

func PrintSomething() {
	fmt.Print("hello world\n")
}

// obstruction is not considered a state, and is handled internally by the door system
func StateKeeper(
	id int,
	initfloor int,
	buttonClick <-chan ButtonEvent,
	floorReached <-chan int,
	motor <-chan PhysicalState,
	mechError <-chan bool,

	ordersWithConsesusToHardware chan<- OrdersWithConsesus,
	physicsToHardware chan<- PhysicalState,

	stateToController chan<- PhysicalState,
	referenceRequest <-chan struct{},
	refToController chan<- PhysicalState,

	netMessageToNetworkSender chan<- NetMessage,
	netMessageToState <-chan NetMessage,
	netErrorToState <-chan NetErrorNotification,
) {
	var wView ElevWorldView = initWorldView(id, initfloor)
	me := &wView.ElevStates[id]
	physicalState := &me.PhysicalState
	stateToController <- *physicalState

	var lastOrdersWithConsensus OrdersWithConsesus
	var lastPhysics PhysicalState
	lastPhysics.Floor = -999

	for {
		// PrintElevState(*me)
		sendToController := true
		select {
		case buttonEvent := <-buttonClick:
			handleButton(&wView, buttonEvent)
		case floorEvent := <-floorReached:
			//fmt.Print("floor update", floorEvent, "\n")
			handleFloor(physicalState, floorEvent)
		case motorEvent := <-motor:
			//fmt.Print("motor update", motorEvent, "\n")
			handleMotor(&wView, motorEvent)
		case mechEvent := <-mechError:
			fmt.Println("mech error", mechEvent)
			handleMech(&wView, mechEvent)
		case netMessage := <-netMessageToState:
			_ = netMessage
		case netErrorNotification := <-netErrorToState: //burde dette caset og det over synkroniseres?
			wView.NetError[netErrorNotification.ID] = netErrorNotification.NetError
			fmt.Println("NetError:", wView.NetError)
		case _ = <-referenceRequest:
			var physics [NumElevators]PhysicalState
			for elev := 0; elev < NumElevators; elev++ {
				physics[elev] = wView.ElevStates[elev].PhysicalState
			}
			ordersWithConsesus := findConsensus(wView)
			if ordersWithConsesus != lastOrdersWithConsensus || me.PhysicalState != lastPhysics {
				relevantOrders := hra.HRA(ordersWithConsesus, physics, wView.NetError)
				ref := referenceGenerator.ReferenceGenerator(me.PhysicalState, relevantOrders)
				_ = ref
				refToController <- ref
			} else {
				sendToController = false
			}
			lastOrdersWithConsensus = ordersWithConsesus
			lastPhysics = me.PhysicalState
		}
		handleOrderDynamics(&wView)

		//Update hardware
		ordersWithConsesus := findConsensus(wView)
		ordersWithConsesusToHardware <- ordersWithConsesus
		physicsToHardware <- *physicalState

		//New state info to network
		var cabBackups [NumElevators][NumFloors]CabOrderState
		for elev := 0; elev < NumElevators; elev++ {
			cabBackups[elev] = wView.ElevStates[elev].OrderState.CabOrders
		}
		netMessage := NetMessage{
			ID:         id,
			ElevState:  *me,
			CabBackups: cabBackups,
		}
		netMessageToNetworkSender <- netMessage

		if sendToController {
			stateToController <- *physicalState
		}

	}
}

func initWorldView(id int, initfloor int) ElevWorldView {
	var wView ElevWorldView

	wView.ID = id
	for elev := 0; elev < NumElevators; elev++ {
		wView.NetError[elev] = true
		wView.CabArchiveSeen[elev] = false

		for floor := 0; floor < NumFloors; floor++ {
			wView.ElevStates[elev].OrderState.HallOrders[floor][Down] = HallNO
			wView.ElevStates[elev].OrderState.HallOrders[floor][Up] = HallNO
			wView.ElevStates[elev].OrderState.CabOrders[floor] = CabUO
		}
	}

	me := &wView.ElevStates[id]
	me.PhysicalState.MechError = false
	me.PhysicalState.Behaviour = Idle
	me.PhysicalState.Floor = initfloor

	return wView
}
