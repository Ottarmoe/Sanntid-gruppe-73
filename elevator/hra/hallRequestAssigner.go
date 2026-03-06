package hra

import (
	"encoding/json"
	"fmt"
	"os/exec"

	. "elevator/elevatorConstants"
	"elevator/state"
	. "elevator/state"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase
type HRAElevstate struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevstate `json:"states"`
}

func HRA(
	orders state.OrdersWithConsesus,
	physics [NumElevators]PhysicalState,
	NetError [3]bool,
) [][2]bool {

	hraExecutable := "hra/hall_request_assigner"
	id := orders.ID

	input := HRAInput{
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		States:       make(map[string]HRAElevstate),
	}

	input.HallRequests = orders.HallOrders[:]

	//
	for elev := 0; elev < NumElevators; elev++ {
		//if no errors, or the elevator is us
		if !physics[elev].MechError && !NetError[elev] || elev == id {

			input.States[fmt.Sprintf("%d", elev)] = HRAElevstate{
				Behavior:    []string{"idle", "moving", "doorOpen"}[physics[id].Behaviour],
				Floor:       physics[id].Floor,
				Direction:   []string{"up", "Down"}[physics[id].MovDirection],
				CabRequests: orders.CabOrders[elev][:],
			}
		}
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
	}

	ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
	}
	//fmt.Printf("output: \n")
	//for k, v := range *output {
	//	fmt.Printf("%6v :  %+v\n", k, v)
	//}

	return (*output)[fmt.Sprintf("%d", id)]
}

func Test() {
	hraExecutable := "hra/hall_request_assigner"

	input := HRAInput{
		HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
		States: map[string]HRAElevstate{
			"one": HRAElevstate{
				Behavior:    "moving",
				Floor:       2,
				Direction:   "up",
				CabRequests: []bool{false, false, false, true},
			},
			"two": HRAElevstate{
				Behavior:    "idle",
				Floor:       0,
				Direction:   "stop",
				CabRequests: []bool{false, false, false, false},
			},
			"three": HRAElevstate{
				Behavior:    "idle",
				Floor:       0,
				Direction:   "stop",
				CabRequests: []bool{false, false, false, false},
			},
		},
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		return
	}

	ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

	// fmt.Printf("output: \n")
	// for k, v := range *output {
	//     fmt.Printf("%6v :  %+v\n", k, v)
	// }
}
