package tests

import (
	// . "elevator/elevatorConstants"
	// "elevator/hallRequestAssigner"
	"elevator/state"
	. "elevator/stateTypes"
	"fmt"
)


func TestSOD() {
	for _, me := range []HallOrderState{HallO, HallOPR, HallNO} {
		for _, a := range []HallOrderState{HallO, HallOPR, HallNO, -999} {
			for _, b := range []HallOrderState{HallO, HallOPR, HallNO, -999} {
				elevatorStates := []HallOrderState{me}
				if a != -999 {
					elevatorStates = append(elevatorStates, a)
				}
				if b != -999 {
					elevatorStates = append(elevatorStates, b)
				}

				fmt.Println(HallOrderStateString(me), HallOrderStateString(a), HallOrderStateString(b), "becomes ", HallOrderStateString(state.SingleOrderDiffusion(me, elevatorStates)))
			}
		}
	}
}

func HallOrderStateString(x HallOrderState) string {
	if x == -999 {
		return ""
	}
	return []string{"HallNO ", "HallO ", "HallOPR "}[x]
}
