package elevatorConstants

var NumFloors int = 4

type SensorEventType int

const (
	HallUpButton SensorEventType = iota
	HallDownButton
	CabButton
	Floor
	MechError
	MotorDir
	MotorBehaviour
)

type SensorEvent struct {
	Eventtype SensorEventType
	Data      int
}
