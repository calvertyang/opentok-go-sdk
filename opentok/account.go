package opentok

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const projectURL = "/v2/project"

type ProjectStatus string

const (
	/**
	 * Set to ACTIVE to use the project API key.
	 */
	ProjectActive ProjectStatus = "ACTIVE"
	/**
	 * Set to SUSPENDED, you will not be able to use the project API key (and any OpenTok sessions created with it).
	 */
	ProjectSuspended ProjectStatus = "SUSPENDED"
)

type Project struct {
	Id                     string `json:"id"`                     // The OpenTok project API key
	UserId                 int    `json:"userId"`                 // The OpenTok account id
	Secret                 string `json:"secret"`                 // The OpenTok project API secret
	Status                 string `json:"status"`                 // Whether the project is active ("VALID", "ACTIVE") or suspended ("SUSPENDED").
	UserStatus             string `json:"userStatus"`             // The OpenTok account status
	Name                   string `json:"name"`                   // The name, if you specified one when creating the project; or an empty string if you did not specify a name
	ContactEmail           string `json:"contactEmail"`           // The OpenTok account email
	CreatedAt              int    `json:"createdAt"`              // The time at which the project was created (a UNIX timestamp, in milliseconds)
	UpdatedAt              int    `json:"updatedAt"`              // The time at which the project was updated (a UNIX timestamp, in milliseconds)
	EnvironmentId          int    `json:"environmentId"`          // The environment id that project is running on
	EnvironmentName        string `json:"environmentName"`        // The environment name that project is running on
	EnvironmentDescription string `json:"environmentDescription"` // The environment description that project is running on
	ApiKey                 string `json:"apiKey"`                 // The OpenTok project API key
}

/**
 * Use this method to create an OpenTok API key and secret for a project.
 */
func (ot *OpenTok) CreateProject(projectName string) (*Project, error) {
	jsonStr := []byte{}
	if projectName != "" {
		jsonStr = []byte(`{ "name": "` + projectName + `" }`)
	}

	//Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	if projectName != "" {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("X-OPENTOK-AUTH", jwt)

	client := &http.Client{}
	res, err := client.Do(req)
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

/**
 * Use this method to get a project details record describing the project (or to get the records for all projects).
 */
func (ot *OpenTok) GetProjectInfo(projectApiKey string) (*[]Project, error) {
	//Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL
	if projectApiKey != "" {
		endpoint += "/" + projectApiKey
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	if projectApiKey != "" {
		// return specific project information
		project := Project{}
		if err = json.NewDecoder(res.Body).Decode(&project); err != nil {
			return nil, err
		}

		projects := []Project{project}

		return &projects, nil
	} else {
		// return all projects information
		projects := []Project{}
		if err = json.NewDecoder(res.Body).Decode(&projects); err != nil {
			return nil, err
		}

		return &projects, nil
	}
}

/**
 * Account administrators can use this method to change a project's status.
 */
func (ot *OpenTok) ChangeProjectStatus(projectApiKey string, projectStatus ProjectStatus) (*Project, error) {
	if projectApiKey == "" {
		return nil, fmt.Errorf("Project status cannot be changed without a project API key")
	}

	if projectStatus == "" || (projectStatus != ProjectActive && projectStatus != ProjectSuspended) {
		return nil, fmt.Errorf("Project status cannot be changed without a valid project status")
	}

	jsonStr := []byte(`{ "status": "` + projectStatus + `" }`)

	//Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL + "/" + projectApiKey
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)

	client := &http.Client{}
	res, err := client.Do(req)
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

/**
 * For security reasons, you may want to generate a new API secret for a project.
 */
func (ot *OpenTok) RefreshProjectSecret(projectApiKey string) (*Project, error) {
	if projectApiKey == "" {
		return nil, fmt.Errorf("Project secret cannot be refreshed without a project API key")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL + "/" + projectApiKey + "/refreshSecret"
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)

	client := &http.Client{}
	res, err := client.Do(req)
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

/**
 * Use this method to delete a project. This prevents the use of the project API key (an any OpenTok sessions created with it).
 */
func (ot *OpenTok) DeleteProject(projectApiKey string) error {
	if projectApiKey == "" {
		return fmt.Errorf("Project cannot be deleted without a project API key")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return err
	}

	endpoint := apiHost + projectURL + "/" + projectApiKey
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	return nil
}
