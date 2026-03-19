package stateTypes

import (
	"fmt"
)

func PrintPhysicalState(stat PhysicalState) {
	switch stat.Behaviour {
	case Idle:
		fmt.Print("Idle ", []string{"Up", "Down"}[stat.MovDirection])
	case Moving:
		fmt.Print("Moving ", []string{"Up", "Down"}[stat.MovDirection])
	case DoorOpen:
		fmt.Print("DoorOpen ", []string{"Up", "Down"}[stat.MovDirection])
	}
	fmt.Println(" on floor", stat.Floor)
}