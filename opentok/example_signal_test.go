package opentok_test

import (
	"fmt"

	"github.com/calvertyang/opentok-go-sdk/opentok"
)

func ExampleOpenTok_SendSessionSignal() {
	err := ot.SendSessionSignal("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", &opentok.SignalData{
		Type: "foo",
		Data: "bar",
	})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleOpenTok_SendConnectionSignal() {
	err := ot.SendConnectionSignal("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", "18145975-97c8-4802-8975-fc8408d56d5c", &opentok.SignalData{
		Type: "foo",
		Data: "bar",
	})
	if err != nil {
		fmt.Println(err)
	}
}
