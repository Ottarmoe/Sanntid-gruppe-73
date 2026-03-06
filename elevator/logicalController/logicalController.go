package logicalController

import (
	. "elevator/elevatorConstants"
	. "elevator/state"
	. "elevator/stateTypes"
	"math"
	"time"
)

func Controller(
	references <-chan PhysicalState,
	actualStates <-chan PhysicalState,
	obstruction <-chan bool,
	referenceRequest chan<- struct{},
	physicalStateUpdate chan<- PhysicalState,
	mechError chan<- bool,
) {

	newDeadLine := make(chan float64)
	go inspector(newDeadLine, mechError)
	doorClosed := make(chan struct{})
	doorOpenFor := make(chan float64)
	go doors(doorOpenFor, doorClosed, obstruction)

	actualState := <-actualStates
	ref := actualState

	stateChanged := false
	for {
		//react and control

		//if we have reached the right floor, but not yet entered the right state
		if actualState.Floor == ref.Floor && actualState.Behaviour != ref.Behaviour {
			if ref.Behaviour == DoorOpen {
				doorOpenFor <- 3.
				actualState.Behaviour = DoorOpen
				stateChanged = true
			} else if ref.Behaviour == Idle && actualState.Behaviour != doorOpen {
				actualState.Behaviour = Idle
				stateChanged = true
			} else if ref.Behaviour == Moving && actualState.Behaviour != doorOpen {
				actualState.Behaviour = Moving
				stateChanged = true
			}
		}
		if ref.MovDirection != actualState.MovDirection {
			ref.MovDirection = actualState.MovDirection
			stateChanged = true
		}
		if ref.Floor != actualState.Floor {
			actualState.Behaviour = Moving
			stateChanged = true
		}

		if stateChanged {
			physicalStateUpdate <- actualState
		}
		stateChanged = false

		if actualState == ref {
			referenceRequest <- struct{}
		}
		//wait for any change in state, or the arrival of a new reference
		if actualState != ref {
			select {
			case newActual := <-actualStates:
				newActual.MechError = false //controller always tries to move as if it is fully functional
				actualState = newActual
			case _ = <-doorClosed:
				actualState.Behaviour = Idle
				stateChanged = true

			case ref := <-references:
				expectedTime := 0.
				expectedTime += math.Abs(float64(ref.Floor-actualState.Floor)) * 7.
				if actualState.Behaviour == DoorOpen {
					expectedTime += 4.
				}
				expectedTime++
				newDeadLine <- expectedTime
			}
		}
	}

}

func burnoutTimer(span float64, burnout chan<- struct{}) {
	time.Sleep(time.Second * time.Duration(span))
	burnout <- struct{}
}

func doors(HoldOpenFor <-chan float64, doorsClosed chan<- struct{}, obstruction <-chan bool) {
	numTimers := 0
	obs := false
	burnoutReturn := make(chan struct{})

	for {
		select {
		case deadline := <-HoldOpenFor:
			go burnoutTimer(deadline, burnoutReturn)
			numTimers++
		case obs := <-obstruction:
			if !obs {
				numTimers++
				go burnoutTimer(3., burnoutReturn)
			}
		case _ = <-burnoutReturn:
			numTimers--
			if numTimers == 0 && !obs {
				doorsClosed <- struct{}
			}
		}
	}
}

func inspector(ExpectedTimeToNextGoal <-chan float64, mechErrorSignal chan<- bool) {
	numTimers := 0
	burnoutReturn := make(chan struct{})

	for {
		select {
		case deadline := <-ExpectedTimeToNextGoal:
			go burnoutTimer(deadline, burnoutReturn)
			if numTimers == 0 {
				mechErrorSignal <- false
			}
			numTimers++
		case _ = <-burnoutReturn:
			numTimers--
			if numTimers == 0 {
				mechErrorSignal <- true
			}
		}
	}
}
