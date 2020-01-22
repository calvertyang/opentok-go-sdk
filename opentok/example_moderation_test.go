package opentok_test

import "fmt"

func ExampleOpenTok_ForceDisconnect() {
	err := ot.ForceDisconnect("40000001", "efdf2fc7-bd6e-4871-9c1d-531f7f6a9486")
	if err != nil {
		fmt.Println(err)
	}
}
