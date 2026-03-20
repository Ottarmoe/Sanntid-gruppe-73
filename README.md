## TODO

Teste network mellom pcer på lab. Teste utdelt packetloss script. 

Endre networklow slik at den er tilpasset laboppsett, og ikke testing i WSL.

Sjekke om elevio trenger å være en egen go module.

Lage debug branch og delivery branch eller lignende funksjonalitet, der blant annet packetloss sim funksjonalitet er fjernet fra networklow, 
samt annen debug funksjonalitet/kommentarer.


Before running the program, make sure to set up an appropriate `elevatorserver`. 
If you get access errors you may need to `chmod -x hallRequestAssigner/hall_request_assigner`


Run the program `go run main.go --id n`, where n is an integer 0, 1 or 2 corresponding to the id of the elevator (0 by default). 

There can at most be 3 elevators at a time, and they must each have a unique id.
The elevators communicate on port 30073, configured in `networkLow.go`