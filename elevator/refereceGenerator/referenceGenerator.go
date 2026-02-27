package referenceGenerator

import (
	"elevator/hra"
	. "elevator/state"
)

func ReferenceGenerator(
	newState <-chan ElevWorldView,
	//motorRef chan<- PhysicalState,
	//inspRef chan<- PhysicalState,
) {
	for {
		wv := <-newState

		//wv.Elevs[0].NetError = false
		//wv.Elevs[1].NetError = false
		//wv.Elevs[2].NetError = false

		HallRequests := hra.HRA(wv)

		_ = HallRequests
		_ = wv

	}
}
