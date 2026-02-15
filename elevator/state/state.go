package state 

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

type SharedState struct {
	id				int
    HallOrders    	[4][2]hallOrderState 
    CabOrders     	[4][3]cabOrderState      
    CabFloor   	 	string    
    CabDir      
	CabBehavoiur
	CabMechError
}
type NetworkElevator struct{
	CabNetError		bool
	//wantBackUp		bool
	State 			SharedState
	//CabBackedUp		[4]bool
}


func finiteStateMachine(
	buttonClick <-chan ButtonEvent, 
	floorReached <-chan int, 
	stopClick <-chan bool, 
	obstructionChange <-chan bool) {
	
	for {
		select {
		case a := <- buttonClick:
			fmt.Printf("%+v\n", a)
			
		case a := <- floorReached:
			fmt.Printf("%+v\n", a)

		case a := <- obstructionChange:
			fmt.Printf("%+v\n", a)
			
		case a := <- stopClick:
			fmt.Printf("%+v\n", a)
		}
    }   
	var localState State 
	var 
}