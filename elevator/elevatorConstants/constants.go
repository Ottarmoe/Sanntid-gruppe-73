package elevatorConstants

import (
	"os"
	"strconv"
)

const NumFloors int = 4
const NumElevators int = 3

var elevatorID int
func ConstantsInit() {
	if len(os.Args) < 2 {
		elevatorID = 0
	} else{
		elevatorID, _ = strconv.Atoi(os.Args[1])
	}
}
func ID() int {
	return elevatorID
}
