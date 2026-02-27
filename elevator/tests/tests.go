package tests

import (
    "os"
    . "elevator/elevatorConstants"
    "elevator/utilities"
    "elevator/hra"
    "elevio"
    "strconv"
    "fmt"
)

func TestMultipleServers(){
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
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)

    //Example usage of button polling
    for {
        select {
        case a := <- drv_buttons:
            // fmt.Printf("%+v\n", a)
            elevio.SetButtonLamp(a.Button, a.Floor, true)

        case a := <- drv_floors:
            fmt.Printf("%+v\n", a)
            if a == NumFloors-1 {
                d = elevio.MD_Down
            } else if a == 0 {
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)

        case a := <- drv_obstr:
            // fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }

        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < NumFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        }
    }
}

func TimeHRA(){
    total, avg := utilities.TimeN(5, hra.Test)
    fmt.Println("Total:", total)
    fmt.Println("Avg:", avg)
}
