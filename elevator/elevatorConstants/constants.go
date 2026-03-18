package elevatorConstants

import (
	"os"
	"strconv"
)

const NumFloors int = 9
const NumElevators int = 3

const BroadcastRate float64 = 0.02          //Send every x seconds
const NetErrorTimerLength float64 = 1       //Need 0 messages within NetErrorTimerLength to mark as netError
const MessagesNeededWithinInterval int = 10 //Need MessagesNeededWithinInterval messages within NetErrorTimerLength to get peer back online
const IntervalsNeeded int = 3
const DoorOpenDuration float64 = 3      //Seconds to hold door open when we arrive at a floor
const DoorObstructionBuffer float64 = 4 //Seconds to wait before marking as mechError after door obstruction is detected
const DeadlineBuffer float64 = 2        // General buffer added to all deadlines
const SecondsPerFloor float64 = 7       // Estimated travel time between adjacent floors

var elevatorID int

func ConstantsInit() {
	if len(os.Args) < 2 {
		elevatorID = 0
	} else {
		elevatorID, _ = strconv.Atoi(os.Args[1])
	}
}
func ID() int {
	return elevatorID
}
