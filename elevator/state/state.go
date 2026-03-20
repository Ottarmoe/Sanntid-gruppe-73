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
	buttonClickCh <-chan ButtonEvent,
	floorReachedCh <-chan int,
	motorStateCh <-chan PhysicalState,
	mechErrorCh <-chan bool,

	ordersWithConsensusToHardwareCh chan<- OrdersWithConsensus,
	physicsToHardwareCh chan<- PhysicalState,

	stateToControllerCh chan<- PhysicalState,
	referenceRequestCh <-chan struct{},
	refToControllerCh chan<- PhysicalState,

	netMessageToNetworkSenderCh chan<- NetMessage,
	netMessageToStateCh <-chan NetMessage,
	netErrorToStateCh <-chan NetErrorNotification,

	stillAliveCh chan<- struct{},
) {
	var wv ElevWorldView = initWorldView(initfloor)

	//heartbeat guarantees state updates at a certain interval for the watchdog timer in main
	heart := time.NewTicker(HeartBeatInerval)

	//initializing communication to controller
	stateToControllerCh <- wv.MyElev().PhysicalState

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
		case buttonEvent := <-buttonClickCh:
			handleButton(&wv, buttonEvent)

		case floorEvent := <-floorReachedCh:
			handleFloor(&wv, floorEvent)

		case motorEvent := <-motorStateCh:
			handleMotor(&wv, motorEvent)

		case mechErrorEvent := <-mechErrorCh:
			handleMech(&wv, mechErrorEvent)

		case netMessage := <-netMessageToStateCh:
			handleNetworkOrders(&wv, netMessage)
			handleNetworkPhysics(&wv, netMessage)

		case netErrorNotification := <-netErrorToStateCh:
			wv.NetError[netErrorNotification.ID] = netErrorNotification.NetError
			fmt.Println("NetError:", wv.NetError)

		case <-referenceRequestCh:
			ordersWithConsensus := findConsensus(&wv)
			relevantOrders := hallRequestAssigner.HRA(ordersWithConsensus, compilePhysicalStates(&wv), wv.NetError)
			ref := referenceGenerator.ReferenceGenerator(wv.MyElev().PhysicalState, relevantOrders)
			//avoid repeat references
			if lastRef != ref {
				refToControllerCh <- ref
			}
			lastRef = ref
			shouldShareNewState = false
		}
		handleOrderDynamics(&wv)
		stillAliveCh <- struct{}{}

		if shouldShareNewState {
			//Update hardware and controller
			ordersWithConsensus := findConsensus(&wv)
			ordersWithConsensusToHardwareCh <- ordersWithConsensus
			physicsToHardwareCh <- wv.MyElev().PhysicalState
			stateToControllerCh <- wv.MyElev().PhysicalState

			//New state info to network
			netMessageToNetworkSenderCh <- wv.CompileNetMessage()
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
