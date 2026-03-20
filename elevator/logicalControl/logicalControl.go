package logicalControl

import (
	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
	"time"
)

func LogicalController(
	referenceCh <-chan PhysicalState,
	stateUpdateCh <-chan PhysicalState,
	obstructionCh <-chan bool,
	referenceRequestOutCh chan<- struct{},
	physicalStateOutCh chan<- PhysicalState,
	mechErrorCh chan<- bool,
) {
	watchDogGoIdleCh := make(chan struct{})
	watchDogNewDeadlineCh := make(chan time.Duration)
	go controlWatchDog(watchDogNewDeadlineCh, watchDogGoIdleCh, mechErrorCh)
	doorClosedEventCh := make(chan struct{})
	openDoorDurationCh := make(chan time.Duration)
	go doors(openDoorDurationCh, doorClosedEventCh, obstructionCh)

	currState := <-stateUpdateCh
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
					openDoorDurationCh <- DoorOpenDuration
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
			physicalStateOutCh <- currState
		}
		if currState == refState && doReferenceRequest {
			referenceRequestOutCh <- struct{}{}
			if currState.Behaviour == Idle {
				watchDogGoIdleCh <- struct{}{}
			}
		}
		//wait for any change in state, or the arrival of a new reference
		initialState = currState
		doReferenceRequest = true
		select {
		case newActual := <-stateUpdateCh:
			newActual.MechError = false //controller always tries to move as if it is fully functional
			currState.Floor = newActual.Floor
		case <-doorClosedEventCh:
			currState.Behaviour = Idle
		case refState = <-referenceCh:
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
				watchDogNewDeadlineCh <- expectedTime
			}
			doReferenceRequest = false
		}
	}

}

func doors(holdOpenForCh <-chan time.Duration, doorsClosedCh chan<- struct{}, obstructionCh <-chan bool) {
	closingTime := time.NewTimer(time.Second)
	closingTime.Stop()
	obstructed := false
	timerActive := false
	for {
		select {
		case obstructed = <-obstructionCh:
			//only start a new closing timer if the previous closingTime has already passed
			if !obstructed && !timerActive {
				closingTime.Reset(PostObstructionOpenTime)
				timerActive = true
			}
		case <-closingTime.C:
			if !obstructed {
				doorsClosedCh <- struct{}{}
			}
			timerActive = false
		case openTime := <-holdOpenForCh:
			closingTime.Stop()
			closingTime.Reset(openTime)
			timerActive = true
		}
	}
}

func controlWatchDog(deadlineCh <-chan time.Duration, goIdleCh <-chan struct{}, mechErrorOutCh chan<- bool) {
	deadLineTimer := time.NewTimer(time.Second)
	deadLineTimer.Stop()
	idle := true

	for {
		select {
		case deadline := <-deadlineCh:
			deadLineTimer.Stop()
			deadLineTimer.Reset(deadline)
			mechErrorOutCh <- false
			idle = false
		case <-deadLineTimer.C:
			if !idle {
				mechErrorOutCh <- true
			}
		case <-goIdleCh:
			idle = true
		}
	}
}
