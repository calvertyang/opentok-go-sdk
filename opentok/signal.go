package opentok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SignalData defines the type and data of signal
type SignalData struct {
	// The type of the signal.
	// This is a string value that clients can filter on when listening for signals
	Type string `json:"type"`

	// The data of the signal
	Data string `json:"data"`
}

// SendSessionSignal send signals to all participants in an active OpenTok session.
func (ot *OpenTok) SendSessionSignal(sessionID string, data *SignalData) error {
	return ot.SendSessionSignalContext(context.Background(), sessionID, data)
}

// SendSessionSignalContext uses ctx for HTTP requests.
func (ot *OpenTok) SendSessionSignalContext(ctx context.Context, sessionID string, data *SignalData) error {
	if sessionID == "" {
		return fmt.Errorf("Signal cannot be sent without a session ID")
	}

	jsonStr, _ := json.Marshal(data)

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/signal"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
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

// SendConnectionSignal send signals to a specific client in an active OpenTok session.
func (ot *OpenTok) SendConnectionSignal(sessionID, connectionID string, data *SignalData) error {
	return ot.SendConnectionSignalContext(context.Background(), sessionID, connectionID, data)
}

// SendConnectionSignalContext uses ctx for HTTP requests.
func (ot *OpenTok) SendConnectionSignalContext(ctx context.Context, sessionID, connectionID string, data *SignalData) error {
	if sessionID == "" {
		return fmt.Errorf("Signal cannot be sent without a session ID")
	}

	if connectionID == "" {
		return fmt.Errorf("Signal cannot be sent without a connection ID")
	}

	jsonStr, _ := json.Marshal(data)

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/connection/" + connectionID + "/signal"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
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
