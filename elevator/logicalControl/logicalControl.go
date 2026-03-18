package logicalControl

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
	"fmt"
	"math"
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
	watchDogNewDeadline := make(chan float64)
	go watchdog(watchDogNewDeadline, watchDogGoIdle, mechError)
	doorClosedEvent := make(chan struct{})
	openDoorDuration := make(chan float64)
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
		}else{
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
				expectedTime := 0.
				expectedTime += math.Abs(float64(referenceState.Floor-currentPhysicalState.Floor)) * SecondsPerFloor
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

func burnoutTimer(span float64, burnout chan<- struct{}) {
	time.Sleep(time.Second * time.Duration(span))
	burnout <- struct{}{}
}

func doors(holdOpenFor <-chan float64, doorsClosed chan<- struct{}, obstruction <-chan bool) {
	numTimers := 0
	obs := false
	burnoutReturn := make(chan struct{})

	for {
		select {
		case deadline := <-holdOpenFor:
			go burnoutTimer(deadline, burnoutReturn)
			numTimers++
		case obs = <-obstruction:
			if !obs {
				numTimers++
				go burnoutTimer(DoorOpenDuration, burnoutReturn)
			}
		case <-burnoutReturn:
			numTimers--
			if numTimers == 0 && !obs {
				doorsClosed <- struct{}{}
			}
		}
	}
}

func watchdog(deadline <-chan float64, goIdle <-chan struct{}, mechErrorSignal chan<- bool) {
	numTimers := 0
	burnoutReturn := make(chan struct{})
	idle := true

	for {
		select {
		case deadline := <-deadline:
			go burnoutTimer(deadline, burnoutReturn)
			if numTimers == 0 {
				mechErrorSignal <- false
			}
			numTimers++
			idle = false
			//fmt.Println("no longer idle")
		case <-burnoutReturn:
			numTimers--
			if numTimers == 0 && !idle {
				mechErrorSignal <- true
			}

		case <-goIdle:
			//fmt.Println("gone idle")
			idle = true
		}
	}
}
