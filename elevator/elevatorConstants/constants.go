package elevatorConstants

import (
	"os"
	"strconv"
)

const NumFloors int = 4
const NumElevators int = 3

const BroadcastRate float64 = 0.02 //Send every x seconds
const NetErrorTimerLength float64 = 1 //Need 0 messages within NetErrorTimerLength to mark as netError
const MessagesNeededWithinInterval int = 10 //Need MessagesNeededWithinInterval messages within NetErrorTimerLength to get peer back online
const IntervalsNeeded int = 3

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
