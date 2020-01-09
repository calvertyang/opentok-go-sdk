package opentok

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SignalData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

/**
 * Send signals to all participants in an active OpenTok session
 */
func (ot *OpenTok) SendSessionSignal(sessionId string, data SignalData) error {
	if sessionId == "" {
		return fmt.Errorf("Signal cannot be sent without a session ID")
	}

	jsonStr, _ := json.Marshal(data)

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return err
	}

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionId + "/signal"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
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

/**
 * Send signals  to a specific client in an active OpenTok session
 */
func (ot *OpenTok) SendConnectionSignal(sessionId, connectionId string, data SignalData) error {
	if sessionId == "" {
		return fmt.Errorf("Signal cannot be sent without a session ID")
	}

	if connectionId == "" {
		return fmt.Errorf("Signal cannot be sent without a connection ID")
	}

	jsonStr, _ := json.Marshal(data)

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return err
	}

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionId + "/connection/" + connectionId + "/signal"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
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
