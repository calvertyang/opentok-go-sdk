package opentok_test

import (
	"fmt"

	"github.com/calvertyang/opentok-go-sdk/v2/opentok"
)

func ExampleOpenTok_ListStreams() {
	streams, err := ot.ListStreams("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", streams)
	}

	// &opentok.StreamList{
	// 	Count:1,
	// 	Items:[]*opentok.Stream{
	// 		&opentok.Stream{
	// 			ID:              "d962b966-964d-4f18-be3f-e0b181a43b0e",
	// 			VideoType:       "camera",
	// 			Name:            "",
	// 			LayoutClassList: []string{}
	// 		}
	// 	}
	// }
}

func ExampleOpenTok_GetStream() {
	stream, err := ot.GetStream("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", "d962b966-964d-4f18-be3f-e0b181a43b0e")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", stream)
	}

	// &opentok.Stream{
	// 	ID:              "d962b966-964d-4f18-be3f-e0b181a43b0e",
	// 	VideoType:       "camera",
	// 	Name:            "",
	// 	LayoutClassList: []string{}
	// }
}

func ExampleOpenTok_SetStreamClassLists() {
	streams, err := ot.SetStreamClassLists("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", &opentok.StreamClassOptions{
		Items: []*opentok.StreamClass{
			{
				ID:              "d962b966-964d-4f18-be3f-e0b181a43b0e",
				LayoutClassList: []string{"full"},
			},
		},
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", streams)
	}

	// &opentok.StreamList{
	// 	Count:1,
	// 	Items:[]*opentok.Stream{
	// 		&opentok.Stream{
	// 			ID:              "d962b966-964d-4f18-be3f-e0b181a43b0e",
	// 			VideoType:       "camera",
	// 			Name:            "",
	// 			LayoutClassList: []string{"full"}
	// 		}
	// 	}
	// }
}
