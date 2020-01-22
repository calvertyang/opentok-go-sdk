package opentok_test

import (
	"fmt"

	"github.com/calvertyang/opentok-go-sdk/opentok"
)

func ExampleOpenTok_StartArchive() {
	archive, err := ot.StartArchive("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", &opentok.ArchiveOptions{
		Name: "example",
		Layout: &opentok.Layout{
			Type: opentok.PIP,
		},
		Resolution: opentok.HD,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", archive)
	}

	// &opentok.Archive{
	// 	CreatedAt:  1579163008000,
	// 	Duration:   0,
	// 	HasAudio:   true,
	// 	HasVideo:   true,
	// 	ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
	// 	Name:       "example",
	// 	OutputMode: "composed",
	// 	ProjectID:  40000001,
	// 	Reason:     "",
	// 	Resolution: "1280x720",
	// 	SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 	Size:       0,
	// 	Status:     "started",
	// 	URL:        nil,
	// }
}

func ExampleOpenTok_StopArchive() {
	archive, err := ot.StopArchive("c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", archive)
	}

	// &opentok.Archive{
	// 	CreatedAt:  1579163008000,
	// 	Duration:   0,
	// 	HasAudio:   true,
	// 	HasVideo:   true,
	// 	ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
	// 	Name:       "example",
	// 	OutputMode: "composed",
	// 	ProjectID:  40000001,
	// 	Reason:     "user initiated",
	// 	Resolution: "1280x720",
	// 	SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 	Size:       0,
	// 	Status:     "stopped",
	// 	URL:        nil,
	// }
}

func ExampleOpenTok_ListArchives() {
	archiveList, err := ot.ListArchives(&opentok.ArchiveListOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", archiveList)
	}

	// &opentok.ArchiveList{
	// 	Count: 1,
	// 	Items: []*opentok.Archive{
	// 		&opentok.Archive{
	// 			CreatedAt:  1579163008000,
	// 			Duration:   34,
	// 			HasAudio:   true,
	// 			HasVideo:   true,
	// 			ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
	// 			Name:       "example",
	// 			OutputMode: "composed",
	// 			ProjectID:  40000001,
	// 			Reason:     "user initiated",
	// 			Resolution: "1280x720",
	// 			SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 			Size:       355816,
	// 			Status:     "uploaded",
	// 			URL:        nil,
	// 		},
	// 	},
	// }
}

func ExampleOpenTok_GetArchive() {
	archive, err := ot.GetArchive("c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", archive)
	}

	// &opentok.Archive{
	// 	CreatedAt:  1579163008000,
	// 	Duration:   34,
	// 	HasAudio:   true,
	// 	HasVideo:   true,
	// 	ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
	// 	Name:       "example",
	// 	OutputMode: "composed",
	// 	ProjectID:  40000001,
	// 	Reason:     "user initiated",
	// 	Resolution: "1280x720",
	// 	SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 	Size:       355816,
	// 	Status:     "uploaded",
	// 	URL:        nil,
	// }
}

func ExampleOpenTok_DeleteArchive() {
	err := ot.DeleteArchive("c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea")
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleOpenTok_SetArchiveStorage() {
	storageConfig, err := ot.SetArchiveStorage(&opentok.StorageOptions{
		Type: "s3",
		Config: &opentok.AmazonS3Config{
			AccessKey: "myUsername",
			SecretKey: "myPassword",
			Bucket:    "example",
		},
		Fallback: "none",
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", storageConfig)
	}

	// &opentok.StorageOptions{
	// 	Type:   "s3",
	// 	Config: map[string]interface {}{
	// 		"accessKey": "myUsername",
	// 		"bucket":    "example",
	// 		"secretKey": "myPassword"
	// 	},
	// 	Fallback: "none",
	// }
}

func ExampleOpenTok_DeleteArchiveStorage() {
	err := ot.DeleteArchiveStorage()

	if err != nil {
		fmt.Println(err)
	}
}

func ExampleOpenTok_SetArchiveLayout() {
	archive, err := ot.SetArchiveLayout("c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea", &opentok.Layout{
		Type: opentok.BestFit,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", archive)
	}

	// &opentok.Archive{
	// 	CreatedAt:  1579163008000,
	// 	Duration:   0,
	// 	HasAudio:   true,
	// 	HasVideo:   true,
	// 	ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
	// 	Name:       "example",
	// 	OutputMode: "composed",
	// 	ProjectID:  40000001,
	// 	Reason:     "",
	// 	Resolution: "1280x720",
	// 	SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
	// 	Size:       0,
	// 	Status:     "started",
	// 	URL:        nil,
	// }
}
