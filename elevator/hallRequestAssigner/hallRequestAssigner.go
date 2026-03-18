package hallRequestAssigner

import (
	"encoding/json"
	"fmt"
	"os/exec"

	. "elevator/elevatorConstants"
	. "elevator/stateTypes"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase
type HRAElevatorState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevatorState `json:"states"`
}

func HRA(orders OrdersWithConsensus, physics [NumElevators]PhysicalState, netError [NumElevators]bool) OurOrders {
    input := buildHRAInput(orders, physics, netError)
    output, err := runHRAExecutable(input)
    if err != nil {
        fmt.Println("HRA error: ", err)
        return OurOrders{}
    }
    return extractHRAOrders(output, orders)
}


func buildHRAInput(orders OrdersWithConsensus, physics [NumElevators]PhysicalState, netError [NumElevators]bool) HRAInput {
    input := HRAInput{
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
        States:       make(map[string]HRAElevatorState),
    }
	input.HallRequests = orders.HallOrders[:]

    for elev := 0; elev < NumElevators; elev++ {
        if !physics[elev].MechError && !netError[elev] || elev == orders.ID {
            input.States[fmt.Sprintf("%d", elev)] = HRAElevatorState{
                Behavior:    []string{"idle", "moving", "doorOpen"}[physics[elev].Behaviour],
                Floor:       physics[elev].Floor,
                Direction:   []string{"up", "down"}[physics[elev].MovDirection],
                CabRequests: orders.CabOrders[elev][:],
            }
        }
    }
    return input
}

// runHRAExecutable marshals input to JSON, calls the external assigner binary, and unmarshals the result.
func runHRAExecutable(input HRAInput) (map[string][][2]bool, error) {
    hraExecutable := "hallRequestAssigner/hall_request_assigner"
	
	jsonBytes, err := json.Marshal(input)
    if err != nil {
        return nil, fmt.Errorf("json.Marshal: %w", err)
    }
    ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("exec: %w — output: %s", err, string(ret))
    }
    var output map[string][][2]bool
    if err := json.Unmarshal(ret, &output); err != nil {
        return nil, fmt.Errorf("json.Unmarshal: %w", err)
    }
    return output, nil
}

// extractHRAOrders picks this elevator's assigned hall orders from the assigner output and combines with cab orders.
func extractHRAOrders(output map[string][][2]bool, orders OrdersWithConsensus) OurOrders {
    var hallOrders [NumFloors][2]bool
    for i := range hallOrders {
        hallOrders[i] = output[fmt.Sprintf("%d", orders.ID)][i]
    }
    return OurOrders{
        HallOrders: hallOrders,
        CabOrders:  orders.CabOrders[orders.ID],
    }
}