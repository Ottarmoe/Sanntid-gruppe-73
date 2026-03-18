package logicalControl

import (
	. "elevator/elevatorConstants"
	//. "elevator/state"
	. "elevator/stateTypes"
	// "fmt"
	"math"
	"time"
)

func Controller(
	referenceCh <-chan PhysicalState,
	stateUpdateCh <-chan PhysicalState,
	obstructionCh <-chan bool,
	referenceRequestOutCh chan<- struct{},
	physicalStateOutCh chan<- PhysicalState,
	mechError chan<- bool,
) {
	watchDogGoIdleCh := make(chan struct{})
	watchDogNewDeadLineCh := make(chan float64)
	go watchdog(watchDogNewDeadLineCh, watchDogGoIdleCh, mechError)
	doorClosedEventCh := make(chan struct{})
	openDoorDurationCh := make(chan float64)
	go doors(openDoorDurationCh, doorClosedEventCh, obstructionCh)

	currentPhysicalState := <-stateUpdateCh
	referenceState := currentPhysicalState
	initialState := currentPhysicalState
	doReferenceRequest := true

	for {
		//if we have reached the right floor
		if currentPhysicalState.Floor == referenceState.Floor {
			currentPhysicalState.MovDirection = referenceState.MovDirection
			if currentPhysicalState.Behaviour != referenceState.Behaviour {
				if referenceState.Behaviour == DoorOpen {
					openDoorDurationCh <- DoorOpenDuration
					currentPhysicalState.Behaviour = DoorOpen
				} else if referenceState.Behaviour == Idle && currentPhysicalState.Behaviour != DoorOpen {
					currentPhysicalState.Behaviour = Idle
				} else if referenceState.Behaviour == Moving && currentPhysicalState.Behaviour != DoorOpen {
					currentPhysicalState.Behaviour = Moving
				}
			}
		}
		//if we are not yet on the right floor
		if referenceState.Floor != currentPhysicalState.Floor {
			if referenceState.Floor < currentPhysicalState.Floor {
				currentPhysicalState.MovDirection = Down
			}
			if referenceState.Floor > currentPhysicalState.Floor {
				currentPhysicalState.MovDirection = Up
			}
			if currentPhysicalState.Behaviour != Moving {
			}
			currentPhysicalState.Behaviour = Moving

		}
		//if anything was changed
		if initialState != currentPhysicalState {
			physicalStateOutCh <- currentPhysicalState
		}
		if currentPhysicalState == referenceState && doReferenceRequest {
			referenceRequestOutCh <- struct{}{}
			if currentPhysicalState.Behaviour == Idle {
				watchDogGoIdleCh <- struct{}{}
			}
		}
		//wait for any change in state, or the arrival of a new reference
		initialState = currentPhysicalState
		doReferenceRequest = true
		select {
		case newActual := <-stateUpdateCh:
			newActual.MechError = false //controller always tries to move as if it is fully functional
			currentPhysicalState.Floor = newActual.Floor
		case <-doorClosedEventCh:
			currentPhysicalState.Behaviour = Idle
		case referenceState = <-referenceCh:
			referenceState.MechError = false
			if referenceState != currentPhysicalState {
				expectedTime := 0.
				expectedTime += math.Abs(float64(referenceState.Floor-currentPhysicalState.Floor)) * SecondsPerFloor
				if currentPhysicalState.Behaviour == DoorOpen {
					expectedTime += DoorObstructionBuffer //adjust this to adjust sensitivity to obstruction
				}
				expectedTime += DeadlineBuffer
				watchDogNewDeadLineCh <- expectedTime
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

func watchdog(deadlineCh <-chan float64, goIdle <-chan struct{}, mechErrorSignal chan<- bool) {
	numTimers := 0
	burnoutReturn := make(chan struct{})
	idle := true

	for {
		select {
		case deadline := <-deadlineCh:
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
