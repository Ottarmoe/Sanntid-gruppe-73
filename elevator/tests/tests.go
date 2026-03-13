package tests

import (
	. "elevator/elevatorConstants"
	"elevator/hra"
	"elevator/state"
	. "elevator/stateTypes"
	"elevator/utilities"
	"elevio"
	"fmt"
	"os"
	"strconv"
)

func TestMultipleServers() {
	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	serverAdress := fmt.Sprintf("localhost:%d", 15657+id)
	elevio.Init(serverAdress, NumFloors)
	// _ = networkLow.Init()
	// network.TestNodeCommunication(id)

	var d elevio.MotorDirection = elevio.MD_Up
	elevio.SetMotorDirection(d)

	//Button polling
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	//Example usage of button polling
	for {
		select {
		case a := <-drv_buttons:
			// fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == NumFloors-1 {
				d = elevio.MD_Down
			} else if a == 0 {
				d = elevio.MD_Up
			}
			elevio.SetMotorDirection(d)

		case a := <-drv_obstr:
			// fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < NumFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}

func TestSOD() {
	for _, me := range []HallOrderState{HallO, HallOPR, HallNO} {
		for _, a := range []HallOrderState{HallO, HallOPR, HallNO, -999} {
			for _, b := range []HallOrderState{HallO, HallOPR, HallNO, -999} {
				elevatorStates := []HallOrderState{me}
				if a != -999 {
					elevatorStates = append(elevatorStates, a)
				}
				if b != -999 {
					elevatorStates = append(elevatorStates, b)
				}

				fmt.Println(HallOrderStateString(me), HallOrderStateString(a), HallOrderStateString(b), "becomes ", HallOrderStateString(state.SingleOrderDiffusion(me, elevatorStates)))
			}
		}
	}
}

func HallOrderStateString(x HallOrderState) string {
	if x == -999 {
		return ""
	}
	return []string{"HallNO ", "HallO ", "HallOPR "}[x]
}

func TimeHRA() {
	total, avg := utilities.TimeN(5, hra.Test)
	fmt.Println("Total:", total)
	fmt.Println("Avg:", avg)
}

// var d elevio.MotorDirection = elevio.MD_Up
// elevio.SetMotorDirection(d)
// sense_buttons1 := make(chan elevio.ButtonEvent)
// sense_floor1 := make(chan int)
// sense_obstr1 := make(chan bool)
// sense_stop1 := make(chan bool)
// go elevio.PollButtons(sense_buttons1)
// go elevio.PollFloorSensor(sense_floor1)
// go elevio.PollObstructionSwitch(sense_obstr1)
// go elevio.PollStopButton(sense_stop1)
// for {
//     select {
//     case a := <- sense_buttons:
//         // fmt.Printf("%+v\n", a)
//         elevio.SetButtonLamp(a.Button, a.Floor, true)

//     case a := <- sense_floor1:
//         // fmt.Printf("%+v\n", a)
//         if a == NumFloors-1 {
//             d = elevio.MD_Down
//         } else if a == 0 {
//             d = elevio.MD_Up
//         }
//         elevio.SetMotorDirection(d)

//     case a := <- sense_obstr1:
//         // fmt.Printf("%+v\n", a)
//         if a {
//             elevio.SetMotorDirection(elevio.MD_Stop)
//         } else {
//             elevio.SetMotorDirection(d)
//         }

//     case _ = <- sense_stop1:
//         // fmt.Printf("%+v\n", a)
//         for f := 0; f < NumFloors; f++ {
//             for b := elevio.ButtonType(0); b < 3; b++ {
//                 elevio.SetButtonLamp(b, f, false)
//             }
//         }
//     }
// }
