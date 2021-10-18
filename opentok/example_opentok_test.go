package opentok_test

import (
	"github.com/calvertyang/opentok-go-sdk/v2/opentok"
)

const (
	apiKey    = "<your api key here>"
	apiSecret = "<your api secret here>"
)

var ot = opentok.New(apiKey, apiSecret)

func ExampleNew() {
	const (
		apiKey    = "12345678"
		apiSecret = "ba7816bf8f01cfea414140de5dae2223b00361a3"
	)

	ot = opentok.New(apiKey, apiSecret)
}
