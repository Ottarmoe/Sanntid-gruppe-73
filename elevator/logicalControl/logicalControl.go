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

	currentPhysicalState := <-stateUpdate
	referenceState := currentPhysicalState
	initialState := currentPhysicalState
	doReferenceRequest := true

	for {
		//if we have reached the right floor
		if currentPhysicalState.Floor == referenceState.Floor {
			currentPhysicalState.MovDirection = referenceState.MovDirection
			if currentPhysicalState.Behaviour != referenceState.Behaviour {
				if referenceState.Behaviour == DoorOpen {
					openDoorDuration <- DoorOpenDuration
					currentPhysicalState.Behaviour = DoorOpen
				} else if referenceState.Behaviour == Idle && currentPhysicalState.Behaviour != DoorOpen {
					currentPhysicalState.Behaviour = Idle
				} else if referenceState.Behaviour == Moving && currentPhysicalState.Behaviour != DoorOpen {
					currentPhysicalState.Behaviour = Moving
				}
			}
			//if we are not yet on the right floor
		} else {
			if referenceState.Floor < currentPhysicalState.Floor {
				currentPhysicalState.MovDirection = Down
			}
			if referenceState.Floor > currentPhysicalState.Floor {
				currentPhysicalState.MovDirection = Up
			}
			if currentPhysicalState.Behaviour != Moving {
				currentPhysicalState.Behaviour = Moving
			}
		}
		//if anything was changed
		if initialState != currentPhysicalState {
			fmt.Print("S ")
			PrintPhysicalState(referenceState)
			physicalStateOut <- currentPhysicalState
		}
		if currentPhysicalState == referenceState && doReferenceRequest {
			referenceRequestOut <- struct{}{}
			if currentPhysicalState.Behaviour == Idle {
				watchDogGoIdle <- struct{}{}
			}
		}
		//wait for any change in state, or the arrival of a new reference
		initialState = currentPhysicalState
		doReferenceRequest = true
		select {
		case newActual := <-stateUpdate:
			newActual.MechError = false //controller always tries to move as if it is fully functional
			currentPhysicalState.Floor = newActual.Floor
		case <-doorClosedEvent:
			currentPhysicalState.Behaviour = Idle
		case referenceState = <-reference:
			referenceState.MechError = false
			if referenceState != currentPhysicalState {
				//time from traversing between floors
				expectedTime := time.Duration(referenceState.Floor-currentPhysicalState.Floor) * SecondsPerFloor
				if expectedTime < 0 {
					expectedTime = -expectedTime
				}
				//time from door open
				if currentPhysicalState.Behaviour == DoorOpen {
					expectedTime += DoorObstructionBuffer //adjust this to adjust sensitivity to obstruction
				}
				expectedTime += DeadlineBuffer
				watchDogNewDeadline <- expectedTime
				fmt.Printf("R ")
				PrintPhysicalState(referenceState)
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
			doorsClosed <- struct{}{}
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
