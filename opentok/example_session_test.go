package opentok_test

import (
	"fmt"
	"time"

	"github.com/MYM-Tech/opentok-go-sdk/v2/opentok"
)

func ExampleOpenTok_CreateSession() {
	session, err := ot.CreateSession(&opentok.SessionOptions{
		ArchiveMode: opentok.AutoArchived,
		MediaMode:   opentok.Routed,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", session)
	}

	// &opentok.Session{
	// 	SessionID:      "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 	ProjectID:      "40000001",
	// 	CreateDt:       "Wed Jan 01 00:00:00 PST 2020",
	// 	MediaServerURL: "",
	// }
}

func ExampleOpenTok_GenerateToken() {
	token, err := ot.GenerateToken("40000001", &opentok.TokenOptions{
		Role:       opentok.Publisher,
		ExpireTime: time.Now().UTC().Add(1 * 24 * time.Hour).Unix(),
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", token)
	}

	// T1==cGFydG5lcl9pZD08eW91ciBhcGkga2V5IGhlcmU+JnNpZz0yYjQyMzlkNjU4YTVkYmE0NGRhMGMyMmUzOTA2MWM5ZWI1ODQ1MTE1OmNvbm5lY3Rpb25fZGF0YT1mb28lM0RiYXImY3JlYXRlX3RpbWU9MTU3Nzg2NTYwMCZleHBpcmVfdGltZT0xNTc3ODY1NjAwJmluaXRpYWxfbGF5b3V0X2NsYXNzX2xpc3Q9Jm5vbmNlPTAuNDk4OTMzNzE3NzEyNjgyMjUmcm9sZT1wdWJsaXNoZXImc2Vzc2lvbl9pZD0xX01YNDBNREF3TURBd01YNS1NVFUzTnpnMk5UWXdNREF3TUg1NE4ySTBPRTFSWjBSbUsxbFJSbkZRVVdnNGRsWm1UMHQtUVg0
}

func ExampleOpenTok_MuteSession() {
	project, err := ot.MuteSession("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", &opentok.MuteOptions{
		Active: true,
		ExcludedStreams: []string{
			"a919b531-bd0e-41fb-8ff0-cdc15684cc93",
			"7f6d8780-741a-4824-98da-16c1f5f1f043",
		},
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", project)
	}

	// &opentok.Project{
	// 	ID:                     "40000001",
	// 	Secret:                 "ba7816bf8f01cfea414140de5dae2223b00361a3",
	// 	Status:                 "VALID",
	// 	Name:                   "example",
	// 	CreatedAt:              1579163008000,
	// 	EnvironmentName:        "default",
	// 	EnvironmentDescription: "Standard Environment",
	// }
}

func ExampleOpenTok_MuteStream() {
	project, err := ot.MuteStream("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", "7f6d8780-741a-4824-98da-16c1f5f1f043")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", project)
	}

	// &opentok.Project{
	// 	ID:                     "40000001",
	// 	Secret:                 "ba7816bf8f01cfea414140de5dae2223b00361a3",
	// 	Status:                 "VALID",
	// 	Name:                   "example",
	// 	CreatedAt:              1579163008000,
	// 	EnvironmentName:        "default",
	// 	EnvironmentDescription: "Standard Environment",
	// }
}
