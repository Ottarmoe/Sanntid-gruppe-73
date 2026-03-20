package logicalControl

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
	"fmt"
	"time"
)

func LogicalController(
	reference <-chan PhysicalState,
	stateUpdate <-chan PhysicalState,
	obstruction <-chan bool,
	referenceRequestOut chan<- struct{},
	physicalStateOut chan<- PhysicalState,
	mechError chan<- bool,
) {
	watchDogGoIdle := make(chan struct{})
	watchDogNewDeadline := make(chan time.Duration)
	go controlWatchDog(watchDogNewDeadline, watchDogGoIdle, mechError)
	doorClosedEvent := make(chan struct{})
	openDoorDuration := make(chan time.Duration)
	go doors(openDoorDuration, doorClosedEvent, obstruction)

	currState := <-stateUpdate
	refState := currState
	initialState := currState
	doReferenceRequest := true

	for {
		//adjust currState toward refState
		if currState.Behaviour == DoorOpen {
			//no action may be taken if doors are open

		} else if currState.Floor == refState.Floor {
			//if we are on the right floor
			currState.MovDirection = refState.MovDirection
			if currState.Behaviour != refState.Behaviour {
				switch refState.Behaviour {
				case DoorOpen:
					openDoorDuration <- DoorOpenDuration
					currState.Behaviour = DoorOpen
				case Idle:
					currState.Behaviour = Idle
				case Moving:
					currState.Behaviour = Moving
				}
			}

		} else {
			//if we are not yet on the right floor
			if refState.Floor < currState.Floor {
				currState.MovDirection = Down
			}
			if refState.Floor > currState.Floor {
				currState.MovDirection = Up
			}
			currState.Behaviour = Moving
		}
		//if anything was changed
		if initialState != currState {
			fmt.Print("S ")
			PrintPhysicalState(refState)
			physicalStateOut <- currState
		}
		if currState == refState && doReferenceRequest {
			referenceRequestOut <- struct{}{}
			if currState.Behaviour == Idle {
				watchDogGoIdle <- struct{}{}
			}
		}
		//wait for any change in state, or the arrival of a new reference
		initialState = currState
		doReferenceRequest = true
		select {
		case newActual := <-stateUpdate:
			newActual.MechError = false //controller always tries to move as if it is fully functional
			currState.Floor = newActual.Floor
		case <-doorClosedEvent:
			currState.Behaviour = Idle
		case refState = <-reference:
			refState.MechError = false
			if refState != currState {
				//time from traversing between floors
				expectedTime := time.Duration(refState.Floor-currState.Floor) * SecondsPerFloor
				if expectedTime < 0 {
					expectedTime = -expectedTime
				}
				//time from door open
				if currState.Behaviour == DoorOpen {
					expectedTime += DoorObstructionBuffer //adjust this to adjust sensitivity to obstruction
				}
				expectedTime += DeadlineBuffer
				watchDogNewDeadline <- expectedTime
				fmt.Printf("R ")
				PrintPhysicalState(refState)
			}
			doReferenceRequest = false
		}
	}

}

func doors(holdOpenFor <-chan time.Duration, doorsClosed chan<- struct{}, obstruction <-chan bool) {
	closingTime := time.NewTimer(time.Second)
	closingTime.Stop()
	obstructed := false
	timerActive := false
	for {
		select {
		case obstructed = <-obstruction:
			//only start a new closing timer if the previous closingTime has already passed
			if !obstructed && !timerActive {
				closingTime.Reset(PostObstructionOpenTime)
				timerActive = true
			}
		case <-closingTime.C:
			if !obstructed {
				doorsClosed <- struct{}{}
			}
			timerActive = false
		case openTime := <-holdOpenFor:
			closingTime.Stop()
			closingTime.Reset(openTime)
			timerActive = true
		}
	}
}

func controlWatchDog(deadline <-chan time.Duration, goIdle <-chan struct{}, mechErrorOut chan<- bool) {
	deadLineTimer := time.NewTimer(time.Second)
	deadLineTimer.Stop()
	idle := true

	for {
		select {
		case deadline := <-deadline:
			deadLineTimer.Stop()
			deadLineTimer.Reset(deadline)
			mechErrorOut <- false
			idle = false
		case <-deadLineTimer.C:
			if !idle {
				mechErrorOut <- true
			}
		case <-goIdle:
			idle = true
		}
	}
}
