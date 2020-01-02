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
	Id          string `json:"id"`          // The OpenTok project API key
	Secret      string `json:"secret"`      // The OpenTok project API secret
	Status      string `json:"status"`      // Whether the project is active ("ACTIVE") or suspended ("SUSPENDED").
	Name        string `json:"name"`        // The name, if you specified one when creating the project; or an empty string if you did not specify a name
	Environment string `json:"environment"` // This is set to "standard" or "enterprise", and it refers to the environment a project is running on.
	CreatedAt   int    `json:"createdAt"`   // The time at which the project was created (a UNIX timestamp, in milliseconds)
}

/**
 * Use this method to create an OpenTok API key and secret for a project.
 */
func (ot *OpenTok) CreateProject(projectName string) (*Project, error) {
	jsonStr := []byte{}
	if projectName != "" {
		jsonStr, _ = json.Marshal(map[string]string{
			"name": projectName,
		})
	}

	//Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s%s", apiHost, projectURL)
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
func (ot *OpenTok) GetProjectInfo(projectApiKey string) (*Project, error) {
	if projectApiKey == "" {
		return nil, fmt.Errorf("Cannot get project information without a project API key")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s%s/%s", apiHost, projectURL, projectApiKey)
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

	project := &Project{}
	if err = json.NewDecoder(res.Body).Decode(project); err != nil {
		return nil, err
	}

	return project, nil
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

	jsonStr, _ := json.Marshal(map[string]string{
		"status": string(projectStatus),
	})

	//Create jwt token
	jwt, err := ot.jwtToken(accountToken)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s%s/%s", apiHost, projectURL, projectApiKey)
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

	endpoint := fmt.Sprintf("%s%s/%s/refreshSecret", apiHost, projectURL, projectApiKey)
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

	endpoint := fmt.Sprintf("%s%s/%s", apiHost, projectURL, projectApiKey)
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
