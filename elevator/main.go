package main

import (
	// "elevator/network"
	// "elevator/networkLow"
	//"elevator/tests"
	. "elevator/elevatorConstants"
	. "elevator/hardwareControl"

	// referenceGenerator "elevator/referenceGenerator"
	"elevator/logicalController"
	"elevator/state"
	"elevio"
	//"elevator/hra"
)

func main() {
	// tests.TestMultipleServers()
	//tests.TimeHRA()
	//serverAdress := fmt.Sprintf("localhost:%d", 15657)
	elevio.Init("localhost:15657", NumFloors)

	sense_buttons := make(chan elevio.ButtonEvent)
	sense_floor := make(chan int)
	sense_obstr := make(chan bool)
	sense_stop := make(chan bool)
	int_mot := make(chan state.PhysicalState)
	int_mech := make(chan bool)

	ref_request := make(chan struct{})
	ref_to_controller := make(chan state.PhysicalState)
	stat_to_controller := make(chan state.PhysicalState)

	ordersWithConsesusToHardware := make(chan state.OrdersWithConsesus)
	physicsToHardware := make(chan state.PhysicalState)

	go elevio.PollButtons(sense_buttons)
	go elevio.PollFloorSensor(sense_floor)
	go elevio.PollObstructionSwitch(sense_obstr)
	go elevio.PollStopButton(sense_stop)

	go state.StateKeeper(0, 0,
		sense_buttons, sense_floor, int_mot, int_mech,
		ordersWithConsesusToHardware, physicsToHardware,
		stat_to_controller, ref_request, ref_to_controller)
	go HardWareControl(physicsToHardware, ordersWithConsesusToHardware)
	go logicalController.Controller(ref_to_controller, stat_to_controller, sense_obstr, ref_request, int_mot, int_mech)
	// go referenceGenerator.ReferenceGenerator(stat_Gen)

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

	select {}
}
