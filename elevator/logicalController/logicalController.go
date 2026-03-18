package logicalController

import (
	//. "elevator/elevatorConstants"
	//. "elevator/state"
	. "elevator/stateTypes"
	// "fmt"
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
	goIdle := make(chan struct{})
	newDeadLine := make(chan float64)
	go inspector(newDeadLine, goIdle, mechError)
	doorClosed := make(chan struct{})
	doorOpenFor := make(chan float64)
	go doors(doorOpenFor, doorClosed, obstruction)

	actualState := <-actualStates
	ref := actualState
	initialState := actualState
	doReferenceRequest := true

	for {
		//react and control

		//if we have reached the right floor, but not yet entered the right state
		if actualState.Floor == ref.Floor {
			actualState.MovDirection = ref.MovDirection
			if actualState.Behaviour != ref.Behaviour {
				if ref.Behaviour == DoorOpen {
					doorOpenFor <- 3.
					actualState.Behaviour = DoorOpen
				} else if ref.Behaviour == Idle && actualState.Behaviour != DoorOpen {
					actualState.Behaviour = Idle
				} else if ref.Behaviour == Moving && actualState.Behaviour != DoorOpen {
					actualState.Behaviour = Moving
				}
			}
		}
		//fmt.Println("i am at", actualState.Floor, "i should be at", ref.Floor)
		if ref.Floor != actualState.Floor {
			if ref.Floor < actualState.Floor {
				actualState.MovDirection = Down
				//fmt.Println("i should move down")
			}
			if ref.Floor > actualState.Floor {
				actualState.MovDirection = Up
				//fmt.Println("i should move up")
			}
			if actualState.Behaviour != Moving {
			}
			actualState.Behaviour = Moving

		}
		if initialState != actualState {
			//fmt.Print("sending new state")
			//PrintPhysicalState(actualState)
			physicalStateUpdate <- actualState
		}
		if actualState == ref && doReferenceRequest {
			// fmt.Println("i have reached my goal")
			referenceRequest <- struct{}{}
			if actualState.Behaviour == Idle {
				goIdle <- struct{}{}
			}
		}
		//wait for any change in state, or the arrival of a new reference
		initialState = actualState
		doReferenceRequest = true
		select {
		case newActual := <-actualStates:
			newActual.MechError = false //controller always tries to move as if it is fully functional
			actualState.Floor = newActual.Floor
			// fmt.Print("S ")
			// PrintPhysicalState(actualState)
			// fmt.Print("r ")
			// PrintPhysicalState(ref)
		case _ = <-doorClosed:
			actualState.Behaviour = Idle
		case ref = <-references:
			ref.MechError = false
			//fmt.Println("is reference", ref, "not the same as actual", actualState)
			if ref != actualState {
				expectedTime := 0.
				expectedTime += math.Abs(float64(ref.Floor-actualState.Floor)) * 7.
				if actualState.Behaviour == DoorOpen {
					expectedTime += 4. //adjust this to adjust sensitivity to obstruction
				}
				expectedTime += 2
				newDeadLine <- expectedTime
				// fmt.Print("R ")
				// PrintPhysicalState(ref)
			}
			doReferenceRequest = false
		}
	}

}

func burnoutTimer(span float64, burnout chan<- struct{}) {
	time.Sleep(time.Second * time.Duration(span))
	burnout <- struct{}{}
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
		case obs = <-obstruction:
			if !obs {
				numTimers++
				go burnoutTimer(3., burnoutReturn)
			}
		case _ = <-burnoutReturn:
			numTimers--
			if numTimers == 0 && !obs {
				doorsClosed <- struct{}{}
			}
		}
	}
}

func inspector(ExpectedTimeToNextGoal <-chan float64, goIdle <-chan struct{}, mechErrorSignal chan<- bool) {
	numTimers := 0
	burnoutReturn := make(chan struct{})
	idle := true

	for {
		select {
		case deadline := <-ExpectedTimeToNextGoal:
			go burnoutTimer(deadline, burnoutReturn)
			if numTimers == 0 {
				mechErrorSignal <- false
			}
			numTimers++
			idle = false
			//fmt.Println("no longer idle")
		case _ = <-burnoutReturn:
			numTimers--
			if numTimers == 0 && !idle {
				mechErrorSignal <- true
			}

		case _ = <-goIdle:
			//fmt.Println("gone idle")
			idle = true
		}
	}
}
