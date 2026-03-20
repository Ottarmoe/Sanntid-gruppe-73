package state

import (
	. "elevator/elevatorConstants"
	"elevator/hallRequestAssigner"
	"elevator/referenceGenerator"
	. "elevator/stateTypes"
	"fmt"
	"time"
	// "time"
)

func PrintSomething() {
	fmt.Print("hello world\n")
}

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
	me := wv.MyState()

	heart := time.NewTicker(HeartBeatInerval)

	//initializing communication to controller
	stateToController <- me.PhysicalState

	var lastRef PhysicalState
	lastRef.Floor = -999

	for {
		stateChanged := true
		select {
		case <-heart.C:
		case buttonEvent := <-buttonClick:
			handleButton(&wv, buttonEvent)
		case floorEvent := <-floorReached:
			handleFloor(&wv, floorEvent)
		case motorEvent := <-motor:
			handleMotor(&wv, motorEvent)
		case mechEvent := <-mechError:
			fmt.Println("mech error", mechEvent)
			handleMech(&wv, mechEvent)
			if mechEvent == true {
				return
			}
		case netMessage := <-netMessageToState:
			handleNetworkOrders(&wv, netMessage)
			handleNetworkPhysics(&wv, netMessage)
		case netErrorNotification := <-netErrorToState: //burde dette caset og det over synkroniseres?
			wv.NetError[netErrorNotification.ID] = netErrorNotification.NetError
			fmt.Println("NetError:", wv.NetError)
		case _ = <-referenceRequest:
			var physics [NumElevators]PhysicalState
			for elev := 0; elev < NumElevators; elev++ {
				physics[elev] = wv.ElevStates[elev].PhysicalState
			}
			ordersWithConsensus := findConsensus(wv)
			relevantOrders := hallRequestAssigner.HRA(ordersWithConsensus, physics, wv.NetError)
			ref := referenceGenerator.ReferenceGenerator(me.PhysicalState, relevantOrders)
			_ = ref
			//fmt.Println("sending ref to controller")
			if lastRef != ref {
				refToController <- ref
			}
			lastRef = ref
			stateChanged = false
		}
		handleOrderDynamics(&wv)
		stillAlive <- struct{}{}

		if stateChanged {
			//Update hardware
			ordersWithConsensus := findConsensus(wv)
			//fmt.Println("sending to hardware")
			ordersWithConsensusToHardware <- ordersWithConsensus
			physicsToHardware <- me.PhysicalState

			//New state info to network
			var cabBackups [NumElevators][NumFloors]CabOrderState
			for elev := 0; elev < NumElevators; elev++ {
				cabBackups[elev] = wv.ElevStates[elev].OrderState.CabOrders
			}
			netMessage := NetMessage{
				ID:         ID(),
				ElevState:  *me,
				CabBackups: cabBackups,
			}
			//fmt.Println("sending to net")
			netMessageToNetworkSender <- netMessage
			//fmt.Println("sending to conntroller")
			stateToController <- me.PhysicalState
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

	me := wv.MyState()
	me.PhysicalState.MechError = false
	me.PhysicalState.Behaviour = Idle
	me.PhysicalState.Floor = initfloor

	return wv
}

// poke main at regular intervals, causing it to send its state to the various modules
// this can help resolve various error states, and enables the watchdog timer in main
func heartBeat(pokeChannel chan<- struct{}, interval time.Duration) {
	for {
		pokeChannel <- struct{}{}
		time.Sleep(interval)
	}
}
