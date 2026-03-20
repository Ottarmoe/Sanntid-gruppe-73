
Before running the program, make sure to set up an appropriate `elevatorserver`. 
If you get access errors you may need to `chmod -x hallRequestAssigner/hall_request_assigner`

Run the program `go run main.go --id n`, where n is an integer 0, 1 or 2 corresponding to the id of the elevator (0 by default). 

There can at most be 3 elevators at a time, and they must each have a unique id.
The elevators communicate on port 30173, configured in `elevatorConstants.go`


The system is structured around the central module `StateKeeper`, which maintains an instance of `ElevWorldView`. It takes in events corresponding with changes in the state of the system, simulates the hall/cab order state machines, and distributes the system state to the physical control modules and network module. 
Additionally, it produces "reference states" for the `logicalControl` module when requested. The `hallRequestAssigner` and `referenceGenerator` modules do not themselves have an internal state, but are called by `stateKeeper` in this process.

`logicalControl` drives the motor and doors of the elevator. It does not do so directly, rather, it sends `physicalState` updates to `stateKeeper`, which are then written to the physical components by the `hardware` module. It drives the motor based on the difference between the recorded physical state and reference state. In normal operation several reference states may need to be achieved for to service any one order. If it fails to do so it produces a MechError (Mechanical Error) message. A new reference cannot be generated before the old one is achieved.

`main` constructs the system of modules and channels that performs the standard behaviours of the program. Additionally it takes on the role of a watchdog timer after initialization is finished. If `stateKeeper` fails to update for more than 2 seconds (configurable in `elevatorConstants`), it kills the program, rather than permiting modules like the network module to keep running independently.

The sending part of the `network` module recieves fully formed `netMessage`s from `stateKeeper`, which it sends at a regular interval. The recieving part decodes recieved `netMessages` and sends them directly to `stateKeeper`. All the order logic is contained in `orderHandling.go`. Incoming netMessages are generally written directly into the appropriate fields in `ElevWorldView` (except when reading the archive of cab orders after reinitialization). The only two functions that change the elevators personal view of the Hall/Cab order are `handleOrderDynamics()` and `handleButton()`.
Additionally, `orderHandling` contains the function `findConsensus()`, which flattens the various order states in `ElevWorldView` into booleans.