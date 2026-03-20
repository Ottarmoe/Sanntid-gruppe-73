package elevatorConstants

import (
	"os"
	"strconv"
	"time"
)

const NumFloors int = 9
const NumElevators int = 3

const BroadcastRate float64 = 0.02                            //Send every x seconds
const NetErrorTimerLength float64 = 1                         //Need 0 messages within NetErrorTimerLength to mark as netError
const MessagesNeededWithinInterval int = 10                   //Need MessagesNeededWithinInterval messages within NetErrorTimerLength to get peer back online
const IntervalsNeeded int = 3                                 //Need IntervalsNeeded intervals with enough messages to get peer back online
const DoorOpenDuration time.Duration = time.Second * 3        //Seconds to hold door open when we arrive at a floor
const DoorObstructionBuffer time.Duration = time.Second * 4   //Seconds to wait before marking as mechError after door obstruction is detected
const DeadlineBuffer time.Duration = time.Second * 2          // General buffer added to all deadlines
const SecondsPerFloor time.Duration = time.Second * 4         // Estimated travel time between adjacent floors
const PostObstructionOpenTime time.Duration = time.Second * 1 //amount of time to wait before closing if obstruction released after a long duration

const HeartBeatInerval time.Duration = 500 * time.Millisecond //the interval that state is guaranteed to distribute messages
const DeathCountDown time.Duration = 2 * time.Second          //interval of no movement from state whereupon it is assumed dead

var elevatorID int

func ConstantsInit() {
	if len(os.Args) < 2 {
		elevatorID = 0
	} else {
		elevatorID, _ = strconv.Atoi(os.Args[1])
	}
}
func MyID() int {
	return elevatorID
}
