package opentok_test

import (
	"fmt"

	"github.com/MYM-Tech/opentok-go-sdk/v2/opentok"
)

func ExampleOpenTok_Dial() {
	sipCall, err := ot.Dial("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", &opentok.DialOptions{
		SIP: &opentok.SIP{
			URI: "sip:user@sip.example.comwhen;transport=tls",
			Headers: &opentok.SIPHeaders{
				"X-Foo": "bar",
			},
			Auth: &opentok.SIPAuth{
				Username: "username",
				Password: "password",
			},
			Secure: true,
		},
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", sipCall)
	}

	// &opentok.SIPCall{
	// 	ID:           "ab31819a-cd52-4da4-8b3e-fb9803337c17",
	// 	ConnectionID: "3a6aa409-bfc5-462c-a9c7-59b72aeebf69",
	// 	StreamID:     "f1a58131-7b2c-4fa8-b2a7-64fdc6b2c0f6",
	// }
}

func ExampleOpenTok_SendDTMF() {
	err := ot.SendDTMF("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", "1713")
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleOpenTok_SendDTMFToClient() {
	err := ot.SendDTMFToClient("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", "3a6aa409-bfc5-462c-a9c7-59b72aeebf69", "1713")
	if err != nil {
		fmt.Println(err)
	}
}
