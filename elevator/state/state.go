package state

import (
	. "elevator/elevatorConstants"
	"elevator/hallRequestAssigner"
	"elevator/referenceGenerator"
	. "elevator/stateTypes"
	"fmt"
	"time"
)

// obstruction is not considered a state, and is handled internally by the door system
func StateKeeper(
	initfloor int,
	buttonClick <-chan ButtonEvent,
	floorReached <-chan int,
	motor <-chan PhysicalState,
	mechError <-chan bool,

	ordersWithConsensusToHardware chan<- OrdersWithConsensus,
	physicsToHardware chan<- PhysicalState,

	stateToController chan<- PhysicalState,
	referenceRequest <-chan struct{},
	refToController chan<- PhysicalState,

	netMessageToNetworkSender chan<- NetMessage,
	netMessageToState <-chan NetMessage,
	netErrorToState <-chan NetErrorNotification,

	stillAlive chan<- struct{},
) {
	var wv ElevWorldView = initWorldView(initfloor)

	//heartbeat guarantees state updates at a certain interval for the watchdog timer in main
	heart := time.NewTicker(HeartBeatInerval)

	//initializing communication to controller
	stateToController <- wv.MyElev().PhysicalState

	//last reference sent, prevents repeat references
	var lastRef PhysicalState
	lastRef.Floor = -999

	for {
		shouldShareNewState := true
		//wait for some manner of event to trigger reprocessing of the state
		//and then send messages with the new state to other modules

		//in order to guarantee that new references for the logical controller are based on the most recent state
		//the StateKeeper is also obligated to service reference requests from the logicalController when relevant
		select {
		case <-heart.C:
		case buttonEvent := <-buttonClick:
			handleButton(&wv, buttonEvent)

		case floorEvent := <-floorReached:
			handleFloor(&wv, floorEvent)

		case motorEvent := <-motor:
			handleMotor(&wv, motorEvent)

		case mechErrorEvent := <-mechError:
			//fmt.Println("mech error", mechErrorEvent)
			handleMech(&wv, mechErrorEvent)

		case netMessage := <-netMessageToState:
			handleNetworkOrders(&wv, netMessage)
			handleNetworkPhysics(&wv, netMessage)

		case netErrorNotification := <-netErrorToState:
			wv.NetError[netErrorNotification.ID] = netErrorNotification.NetError
			fmt.Println("NetError:", wv.NetError)

		case <-referenceRequest:
			ordersWithConsensus := findConsensus(&wv)
			relevantOrders := hallRequestAssigner.HRA(ordersWithConsensus, compilePhysicalStates(&wv), wv.NetError)
			ref := referenceGenerator.ReferenceGenerator(wv.MyElev().PhysicalState, relevantOrders)
			//avoid repeat references
			if lastRef != ref {
				refToController <- ref
			}
			lastRef = ref
			shouldShareNewState = false
		}
		handleOrderDynamics(&wv)
		stillAlive <- struct{}{}

		if shouldShareNewState {
			//Update hardware and controller
			ordersWithConsensus := findConsensus(&wv)
			ordersWithConsensusToHardware <- ordersWithConsensus
			physicsToHardware <- wv.MyElev().PhysicalState
			stateToController <- wv.MyElev().PhysicalState

			//New state info to network
			netMessageToNetworkSender <- wv.CompileNetMessage()
		}
	}
}

func initWorldView(initfloor int) ElevWorldView {
	var wv ElevWorldView

	for elev := 0; elev < NumElevators; elev++ {
		wv.NetError[elev] = true
		wv.CabArchiveSeen[elev] = false

		for floor := 0; floor < NumFloors; floor++ {
			wv.ElevStates[elev].OrderState.HallOrders[floor][Down] = HallNO
			wv.ElevStates[elev].OrderState.HallOrders[floor][Up] = HallNO
			wv.ElevStates[elev].OrderState.CabOrders[floor] = CabUO
		}
	}

	me := wv.MyElev()
	me.PhysicalState.MechError = false
	me.PhysicalState.Behaviour = Idle
	me.PhysicalState.Floor = initfloor

	return wv
}
