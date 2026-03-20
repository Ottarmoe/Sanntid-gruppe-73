package elevatorConstants

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// State
const NumFloors int = 4
const NumElevators int = 3
const HeartBeatInerval time.Duration = 500 * time.Millisecond //the interval that state is guaranteed to distribute messages
const DeathCountDown time.Duration = 2 * time.Second          //interval of no movement from state whereupon it is assumed dead

// Network
const BroadcastRate time.Duration = time.Millisecond * 20
const NetErrorTimerLength time.Duration = time.Second * 1 //Need 0 messages within NetErrorTimerLength to mark as netError
const MessagesNeededWithinInterval int = 10               //Need MessagesNeededWithinInterval messages within NetErrorTimerLength to get peer back online
const IntervalsNeeded int = 3                             //Need IntervalsNeeded intervals with enough messages to get peer back online

// logicalControl
const DoorOpenDuration time.Duration = time.Second * 3        //Seconds to hold door open when we arrive at a floor
const DoorObstructionBuffer time.Duration = time.Second * 1   // acceptable time for door to be open beyond DoorOpenDuration before marking MechError
const DeadlineBuffer time.Duration = time.Second * 2          // General buffer added to all deadlines
const SecondsPerFloor time.Duration = time.Second * 4         // Estimated travel time between adjacent floors
const PostObstructionOpenTime time.Duration = time.Second * 1 //amount of time to wait before closing if obstruction released after a long duration

// Elevator specific configuration, personal id and if it should run on a simulator
var elevatorID int
var usingSimulator bool

func ConstantsInit() {
	flag.IntVar(&elevatorID, "id", 0, "ID of this elevator")
	flag.BoolVar(&usingSimulator, "sim", false, "Run in simulation mode (true/false)")

	flag.Parse()

	// Validate ID
	if elevatorID < 0 || elevatorID >= NumElevators {
		fmt.Printf("Error: --id must be between 0 and %d (got %d)\n", NumElevators-1, elevatorID)
		os.Exit(1)
	}
}
func MyID() int {
	return elevatorID
}
func UsingSimulator() bool {
	return usingSimulator
}
