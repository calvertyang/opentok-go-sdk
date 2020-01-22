package opentok_test

import (
	"fmt"

	"github.com/calvertyang/opentok-go-sdk/v2/opentok"
)

func ExampleOpenTok_StartBroadcast() {
	broadcast, err := ot.StartBroadcast("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", &opentok.BroadcastOptions{
		Layout: &opentok.Layout{
			Type: opentok.PIP,
		},
		MaxDuration: 60,
		Outputs: &opentok.BroadcastOutputOptions{
			RTMP: []*opentok.RTMPConfig{
				&opentok.RTMPConfig{
					ID:         "foo",
					ServerURL:  "rtmps://myfooserver/myfooapp",
					StreamName: "myfoostream",
				},
				&opentok.RTMPConfig{
					ID:         "bar",
					ServerURL:  "rtmp://mybarserver/mybarapp",
					StreamName: "mybarstream",
				},
			},
		},
		Resolution: opentok.HD,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", broadcast)
	}

	// &opentok.Broadcast{
	// 	ID:            "ce872e0d-4997-440a-a0a5-10ce715b54cf",
	// 	SessionID:     "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 	ProjectID:     "40000001",
	// 	CreatedAt:     1579163008000,
	// 	UpdatedAt:     1579163008000,
	// 	Resolution:    "1280x720",
	// 	Status:        "started",
	// 	BroadcastURLs: &opentok.BroadcastURLs{
	// 		HLS: "",
	// 		RTMP:[]*opentok.RTMPConfig{
	// 			&opentok.RTMPConfig{
	// 				ID:         "foo",
	// 				Status:     "connecting",
	// 				ServerURL:  "rtmps://myfooserver/myfooapp",
	// 				StreamName: "myfoostream"
	// 			},
	// 			&opentok.RTMPConfig{
	// 				ID:         "bar",
	// 				Status:     "connecting",
	// 				ServerURL:  "rtmp://mybarserver/mybarapp",
	// 				StreamName: "mybarstream"
	// 			},
	// 		},
	// 	},
	// }
}

func ExampleOpenTok_StopBroadcast() {
	broadcast, err := ot.StopBroadcast("ce872e0d-4997-440a-a0a5-10ce715b54cf")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", broadcast)
	}

	// &opentok.Broadcast{
	// 	ID:            "ce872e0d-4997-440a-a0a5-10ce715b54cf",
	// 	SessionID:     "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 	ProjectID:     40000001,
	// 	CreatedAt:     1579163008000,
	// 	UpdatedAt:     1579163009000,
	// 	Resolution:    "1280x720",
	// 	Status:        "stopped",
	// 	BroadcastURLs: nil,
	// }
}

func ExampleOpenTok_ListBroadcasts() {
	broadcasts, err := ot.ListBroadcasts(&opentok.BroadcastListOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", broadcasts)
	}

	// &opentok.BroadcastList{
	// 	Count:1,
	// 	Items:[]*opentok.Broadcast{
	// 		&opentok.Broadcast{
	// 			ID:            "ce872e0d-4997-440a-a0a5-10ce715b54cf",
	// 			SessionID:     "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 			ProjectID:     "40000001",
	// 			CreatedAt:     1579163008000,
	// 			UpdatedAt:     1579163008000,
	// 			Resolution:    "1280x720",
	// 			Status:        "started",
	// 			BroadcastURLs: &opentok.BroadcastURLs{
	// 				HLS: "",
	// 				RTMP:[]*opentok.RTMPConfig{
	// 					&opentok.RTMPConfig{
	// 						ID:         "foo",
	// 						Status:     "connecting",
	// 						ServerURL:  "rtmps://myfooserver/myfooapp",
	// 						StreamName: "myfoostream"
	// 					},
	// 					&opentok.RTMPConfig{
	// 						ID:         "bar",
	// 						Status:     "connecting",
	// 						ServerURL:  "rtmp://mybarserver/mybarapp",
	// 						StreamName: "mybarstream"
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
}

func ExampleOpenTok_GetBroadcast() {
	broadcast, err := ot.GetBroadcast("ce872e0d-4997-440a-a0a5-10ce715b54cf")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", broadcast)
	}

	// &opentok.Broadcast{
	// 	ID:            "ce872e0d-4997-440a-a0a5-10ce715b54cf",
	// 	SessionID:     "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 	ProjectID:     40000001,
	// 	CreatedAt:     1579163008000,
	// 	UpdatedAt:     1579163009000,
	// 	Resolution:    "1280x720",
	// 	Status:        "stopped",
	// 	BroadcastURLs: nil,
	// }
}

func ExampleOpenTok_SetBroadcastLayout() {
	broadcast, err := ot.SetBroadcastLayout("ce872e0d-4997-440a-a0a5-10ce715b54cf", &opentok.Layout{
		Type: opentok.BestFit,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", broadcast)
	}

	// &opentok.Broadcast{
	// 	ID:            "ce872e0d-4997-440a-a0a5-10ce715b54cf",
	// 	SessionID:     "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 	ProjectID:     "40000001",
	// 	CreatedAt:     1579163008000,
	// 	UpdatedAt:     1579163008000,
	// 	Resolution:    "1280x720",
	// 	Status:        "started",
	// 	BroadcastURLs: &opentok.BroadcastURLs{
	// 		HLS: "",
	// 		RTMP:[]*opentok.RTMPConfig{
	// 			&opentok.RTMPConfig{
	// 				ID:         "foo",
	// 				Status:     "connecting",
	// 				ServerURL:  "rtmps://myfooserver/myfooapp",
	// 				StreamName: "myfoostream"
	// 			},
	// 			&opentok.RTMPConfig{
	// 				ID:         "bar",
	// 				Status:     "connecting",
	// 				ServerURL:  "rtmp://mybarserver/mybarapp",
	// 				StreamName: "mybarstream"
	// 			},
	// 		},
	// 	},
	// }
}
