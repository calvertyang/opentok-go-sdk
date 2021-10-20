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

	// Create jwt token
	jwt, err := ot.genProjectJWT()
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/signal"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)

	res, err := ot.sendRequest(req, ctx)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return parseErrorResponse(res)
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

	// Create jwt token
	jwt, err := ot.genProjectJWT()
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/connection/" + connectionID + "/signal"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)

	res, err := ot.sendRequest(req, ctx)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return parseErrorResponse(res)
	}

	return nil
}

// SendSignal send signals to all participants.
func (s *Session) SendSignal(data *SignalData) error {
	return s.SendSignalContext(context.Background(), data)
}

// SendSignalContext uses ctx for HTTP requests.
func (s *Session) SendSignalContext(ctx context.Context, data *SignalData) error {
	return s.OpenTok.SendSessionSignalContext(ctx, s.SessionID, data)
}

// SendConnectionSignal send signals to a specific client.
func (s *Session) SendConnectionSignal(connectionID string, data *SignalData) error {
	return s.SendConnectionSignalContext(context.Background(), connectionID, data)
}

// SendConnectionSignalContext uses ctx for HTTP requests.
func (s *Session) SendConnectionSignalContext(ctx context.Context, connectionID string, data *SignalData) error {
	return s.OpenTok.SendConnectionSignalContext(ctx, s.SessionID, connectionID, data)
}
