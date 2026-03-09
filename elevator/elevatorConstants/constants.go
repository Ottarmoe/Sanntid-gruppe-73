package elevatorConstants

import (
	"os"
	"strconv"
)

const NumFloors int = 4
const NumElevators int = 3
const BroadcastRate float64 = 0.02 //Send every 20 ms
const NetErrorTimerLength float64 = 1 //1 second
const MessagesNeededWithinInterval int = 10 //Need 10 messages within to get back online

var elevatorID int
func ConstantsInit() {
	if len(os.Args) < 2 {
		elevatorID = 0
	} else{
		elevatorID, _ = strconv.Atoi(os.Args[1])
	}
}
func ID() int {
	return elevatorID
}
