package main

// import . "elevator/elevatorConstants"
// import "elevio"
// import "fmt"
// import "elevator/hra"
// import "elevator/utilities"

// func main(){
//     //Init functions
//     elevio.Init("localhost:15657", NumFloors)
//     var d elevio.MotorDirection = elevio.MD_Stop//elevio.MD_Up
//     elevio.SetMotorDirection(d)
    
//     //Button polling
//     drv_buttons := make(chan elevio.ButtonEvent)
//     drv_floors  := make(chan int)
//     drv_obstr   := make(chan bool)
//     drv_stop    := make(chan bool)    
    
//     go elevio.PollButtons(drv_buttons)
//     go elevio.PollFloorSensor(drv_floors)
//     go elevio.PollObstructionSwitch(drv_obstr)
//     go elevio.PollStopButton(drv_stop)
    
//     //Example usage of hall request assigner
//     total, avg := utilities.TimeN(100, hra.Test)
//     fmt.Println("Total:", total)
//     fmt.Println("Avg:", avg)

//     //Example usage of button polling
//     for {
//         select {
//         case a := <- drv_buttons:
//             // fmt.Printf("%+v\n", a)
//             elevio.SetButtonLamp(a.Button, a.Floor, true)
            
//         case a := <- drv_floors:
//             fmt.Printf("%+v\n", a)
//             if a == NumFloors-1 {
//                 d = elevio.MD_Down
//             } else if a == 0 {
//                 d = elevio.MD_Up
//             }
//             elevio.SetMotorDirection(d)
            
            
//         case a := <- drv_obstr:
//             // fmt.Printf("%+v\n", a)
//             if a {
//                 elevio.SetMotorDirection(elevio.MD_Stop)
//             } else {
//                 elevio.SetMotorDirection(d)
//             }
            
//         case a := <- drv_stop:
//             fmt.Printf("%+v\n", a)
//             for f := 0; f < NumFloors; f++ {
//                 for b := elevio.ButtonType(0); b < 3; b++ {
//                     elevio.SetButtonLamp(b, f, false)
//                 }
//             }
//         }
//     }   
// }

// import "elevator/hra"
// import "elevator/utilities"
// import "fmt"

// func main(){
//     total, avg := utilities.TimeN(100, hra.Test)
//     fmt.Println("Total:", total)
//     fmt.Println("Avg:", avg)
// }

import "elevator/networkLow"
import "fmt"

func main(){
    fmt.Printf("Program started\n")
    _ = networkLow.Init()
    msg := []byte(fmt.Sprintf("Hello"))
    _ = networkLow.Send(msg)
    msg = []byte(fmt.Sprintf("Bye"))
    _ = networkLow.Send(msg)
    _ = networkLow.Send(msg)

    // buf := make([]byte, 2048)
    // n, addr, _ := networkLow.Receive(buf)
    // networkLow.PrintMessage(buf, n, addr)
    // n, addr, _ = networkLow.Receive(buf)
    // networkLow.PrintMessage(buf, n, addr)


}