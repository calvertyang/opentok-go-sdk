package opentok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const projectURL = "/v2/project"

// ProjectStatus is the alias of string type
type ProjectStatus string

const (
	// ProjectActive is used to set the project API key to ACTIVE.
	ProjectActive ProjectStatus = "ACTIVE"

	// ProjectSuspended is used to set the project API key to SUSPENDED, you will
	// not be able to use the project API key (and any OpenTok sessions created
	// with it).
	ProjectSuspended ProjectStatus = "SUSPENDED"
)

// Project defines the response returned from API
type Project struct {
	// The OpenTok project API key
	ID string `json:"id"`

	// The OpenTok project API secret
	Secret string `json:"secret"`

	// Whether the project is active ("VALID", "ACTIVE") or suspended ("SUSPENDED").
	Status string `json:"status"`

	// The name, if you specified one when creating the project; or an empty
	// string if you did not specify a name
	Name string `json:"name"`

	// The time at which the project was created (a UNIX timestamp, in milliseconds)
	CreatedAt int `json:"createdAt"`

	// The environment name that project is running on
	EnvironmentName string `json:"environmentName"`

	// The environment description that project is running on
	EnvironmentDescription string `json:"environmentDescription"`
}

// CreateProject creates an OpenTok API key and secret for a project..
func (ot *OpenTok) CreateProject(projectName string) (*Project, error) {
	return ot.CreateProjectContext(context.Background(), projectName)
}

// CreateProjectContext uses ctx for HTTP requests.
func (ot *OpenTok) CreateProjectContext(ctx context.Context, projectName string) (*Project, error) {
	jsonStr := []byte{}
	if projectName != "" {
		jsonStr = []byte(`{ "name": "` + projectName + `" }`)
	}

	// Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	if projectName != "" {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	project := &Project{}
	if err = json.NewDecoder(res.Body).Decode(project); err != nil {
		return nil, err
	}

	return project, nil
}

// ListProjectsContext uses ctx for HTTP requests.
func (ot *OpenTok) ListProjects() ([]*Project, error) {
	return ot.ListProjectsContext(context.Background())
}

// ListProjectsContext uses ctx for HTTP requests..
func (ot *OpenTok) ListProjectsContext(ctx context.Context) ([]*Project, error) {
	// Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	projects := []*Project{}
	if err = json.NewDecoder(res.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

// GetProject returns a project details record describing the project.
func (ot *OpenTok) GetProject(projectAPIKey string) (*Project, error) {
	return ot.GetProjectContext(context.Background(), projectAPIKey)
}

// GetProjectContext uses ctx for HTTP requests.
func (ot *OpenTok) GetProjectContext(ctx context.Context, projectAPIKey string) (*Project, error) {
	if projectAPIKey == "" {
		return nil, fmt.Errorf("Cannot get project information without a project API key")
	}

	// Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + projectAPIKey
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	project := &Project{}
	if err = json.NewDecoder(res.Body).Decode(project); err != nil {
		return nil, err
	}

	return project, nil
}

// ChangeProjectStatus changes the status of project. The status is
// either active or suspended.
func (ot *OpenTok) ChangeProjectStatus(projectAPIKey string, projectStatus ProjectStatus) (*Project, error) {
	return ot.ChangeProjectStatusContext(context.Background(), projectAPIKey, projectStatus)
}

// ChangeProjectStatusContext uses ctx for HTTP requests.
func (ot *OpenTok) ChangeProjectStatusContext(ctx context.Context, projectAPIKey string, projectStatus ProjectStatus) (*Project, error) {
	if projectAPIKey == "" {
		return nil, fmt.Errorf("Project status cannot be changed without a project API key")
	}

	if projectStatus != ProjectActive && projectStatus != ProjectSuspended {
		return nil, fmt.Errorf("Project status cannot be changed without a valid project status")
	}

	jsonStr := []byte(`{ "status": "` + projectStatus + `" }`)

	// Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + projectAPIKey
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	project := &Project{}
	if err = json.NewDecoder(res.Body).Decode(project); err != nil {
		return nil, err
	}

	return project, nil
}

// RefreshProjectSecret generates a new API secret for a project.
func (ot *OpenTok) RefreshProjectSecret(projectAPIKey string) (*Project, error) {
	return ot.RefreshProjectSecretContext(context.Background(), projectAPIKey)
}

// RefreshProjectSecretContext uses ctx for HTTP requests.
func (ot *OpenTok) RefreshProjectSecretContext(ctx context.Context, projectAPIKey string) (*Project, error) {
	if projectAPIKey == "" {
		return nil, fmt.Errorf("Project secret cannot be refreshed without a project API key")
	}

	// Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + projectAPIKey + "/refreshSecret"
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	project := &Project{}
	if err = json.NewDecoder(res.Body).Decode(project); err != nil {
		return nil, err
	}

	return project, nil
}

// DeleteProjectContext prevents the use of the project API key (and
// any OpenTok sessions created with it).
func (ot *OpenTok) DeleteProject(projectAPIKey string) error {
	return ot.DeleteProjectContext(context.Background(), projectAPIKey)
}

// DeleteProjectContext uses ctx for HTTP requests.
func (ot *OpenTok) DeleteProjectContext(ctx context.Context, projectAPIKey string) error {
	if projectAPIKey == "" {
		return fmt.Errorf("Project cannot be deleted without a project API key")
	}

	// Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + projectAPIKey
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	return nil
}
