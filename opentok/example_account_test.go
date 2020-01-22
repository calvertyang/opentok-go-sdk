package opentok_test

import (
	"fmt"

	"github.com/calvertyang/opentok-go-sdk/opentok"
)

func ExampleOpenTok_CreateProject() {
	project, err := ot.CreateProject("example")
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

func ExampleOpenTok_ListProjects() {
	projects, err := ot.ListProjects()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", projects)
	}

	// [
	// 	&opentok.Project{
	// 		ID:                     "40000001",
	// 		Secret:                 "ba7816bf8f01cfea414140de5dae2223b00361a3",
	// 		Status:                 "VALID",
	// 		Name:                   "example",
	// 		CreatedAt:              1579163008000,
	// 		EnvironmentName:        "default",
	// 		EnvironmentDescription: "Standard Environment",
	// 	},
	// 	&opentok.Project{
	// 		ID:                     "40000002",
	// 		Secret:                 "fcde2b2edba56bf408601fb721fe9b5c338d10ee",
	// 		Status:                 "VALID",
	// 		Name:                   "foo",
	// 		CreatedAt:              1579163506000,
	// 		EnvironmentName:        "default",
	// 		EnvironmentDescription: "Standard Environment",
	// 	},
	// ]
}

func ExampleOpenTok_GetProject() {
	project, err := ot.GetProject("40000001")
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

// Active a project API key
func ExampleOpenTok_ChangeProjectStatus_active() {
	project, err := ot.ChangeProjectStatus("40000001", opentok.ProjectActive)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", project)
	}

	// &opentok.Project{
	// 	ID:                     "40000001",
	// 	Secret:                 "ba7816bf8f01cfea414140de5dae2223b00361a3",
	// 	Status:                 "ACTIVE",
	// 	Name:                   "example",
	// 	CreatedAt:              1579163008000,
	// 	EnvironmentName:        "default",
	// 	EnvironmentDescription: "Standard Environment",
	// }
}

// Suspend a project API key
func ExampleOpenTok_ChangeProjectStatus_suspend() {
	project, err := ot.ChangeProjectStatus("40000001", opentok.ProjectSuspended)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", project)
	}

	// &opentok.Project{
	// 	ID:                     "40000001",
	// 	Secret:                 "ba7816bf8f01cfea414140de5dae2223b00361a3",
	// 	Status:                 "SUSPENDED",
	// 	Name:                   "example",
	// 	CreatedAt:              1579163008000,
	// 	EnvironmentName:        "default",
	// 	EnvironmentDescription: "Standard Environment",
	// }
}

func ExampleOpenTok_RefreshProjectSecret() {
	project, err := ot.RefreshProjectSecret("40000001")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v", project)
	}

	// &opentok.Project{
	// 	ID:                     "40000001",
	// 	Secret:                 "6a1e361fc002c0b2b51a821d7ce11f96e1887aea",
	// 	Status:                 "ACTIVE",
	// 	Name:                   "example",
	// 	CreatedAt:              1579163008000,
	// 	EnvironmentName:        "default",
	// 	EnvironmentDescription: "Standard Environment",
	// }
}

func ExampleOpenTok_DeleteProject() {
	err := ot.DeleteProject("40000001")
	if err != nil {
		fmt.Println(err)
	}
}
