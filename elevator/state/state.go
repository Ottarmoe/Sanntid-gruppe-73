package state 

import . "elevio"

type hallOrderState int

const (
	HallO hallOrderState = iota
	HallNO
	HallOPR
)

type cabOrderState int

const (
	CabO cabOrderState = iota
	CabNO
	CabUO //unknown order
)

type Direction int

const (
	Down Direction = iota
	Up
)

type MotorBehaviour int

const(
	Moving MotorBehaviour = iota
	Idle
	DoorOpen
)

type MotorState struct{
	Behaviour MotorBehaviour
	Direction MotorDirection
}

type ElevState struct {
	ID				int
	NetError		bool
	CabAgreement	[4]bool
    HallOrders    	[4][2]hallOrderState //0 is down, 1 is up, use "direction"
    CabOrders     	[4]cabOrderState      
    CabFloor   	 	int    
    CabDir      	Direction
	CabBehaviour	MotorBehaviour
	CabMechError	bool
}



//obstruction is not considered a state, and is handled internally by the door system
func finiteStateMachine(
	id int,
	initfloor int,
	buttonClick <-chan ButtonEvent, 
	floorReached <-chan int, 
	Motor <-chan MotorState,
	MechError <-chan bool
	//stopClick <-chan bool, 
	//obstructionChange <-chan bool
	) {

	elevs := make([]ElevState, 1) //zeroth element is always us
	elevs[0].ID = id
	elevs[0].NetError = true //trust me bro

	elevs[0].CabMechError = false
	elevs[0].CabBehaviour = Idle
	elevs[0].CabFloor = initfloor
	
	for floor := 0; floor < 4; floor++{
		elevs[0].HallOrders[floor][Down] = HallNO
		elevs[0].HallOrders[floor][Up] = HallNO
		elevs[0].CabOrders[floor] = CabUO
	} 
	
	
	for {
		select {
		case buttonEvent := <- buttonClick:
			handleButton(elevs, buttonEvent)
			
		case floorEvent := <- floorReached:
			handleFloor(elevs, floorEvent)
		case

		}
		// case a := <- obstructionChange:
		// 	fmt.Printf("%+v\n", a)

    }   
}

func handleButton(elevs []ElevState, event ButtonEvent){
	switch event.Button{
	case BT_HallUp:
		if elevs[0].HallOrders[event.Floor][Up] == HallNO{
			elevs[0].HallOrders[event.Floor][Up] = HallO
		}
	case BT_HallDown:
		if elevs[0].HallOrders[event.Floor][Down] == HallNO{
			elevs[0].HallOrders[event.Floor][Down] = HallO
		}
	case BT_Cab:
		elevs[0].CabOrders[event.Floor] = CabO
		for _, elev := range elevs{
			elev.CabAgreement[event.Floor] = false
		}
	}
}

func handleFloor(elevs []ElevState, event int){
	elevs[0].CabFloor = event
}