package sensors

import (
	. "elevator/elevatorConstants"
	"elevio"
)

func SensorLoop(out chan SensorEvent) {
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)

	for {
		select {
		case a := <-drv_buttons:
			if a.Button == elevio.BT_HallUp {
				out <- SensorEvent{Eventtype: HallUpButton, Data: a.Floor}
			} else if a.Button == elevio.BT_HallDown {
				out <- SensorEvent{Eventtype: HallDownButton, Data: a.Floor}
			} else if a.Button == elevio.BT_Cab {
				out <- SensorEvent{Eventtype: CabButton, Data: a.Floor}
			}

		case a := <-drv_floors:
			out <- SensorEvent{Eventtype: Floor, Data: a}
		}
	}
}
