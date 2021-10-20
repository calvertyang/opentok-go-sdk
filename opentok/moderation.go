package opentok

import (
	"context"
	"fmt"
	"net/http"
)

// ForceDisconnect disconnects a client from an OpenTok session via server-side.
func (ot *OpenTok) ForceDisconnect(sessionID, connectionID string) error {
	return ot.ForceDisconnectContext(context.Background(), sessionID, connectionID)
}

// ForceDisconnectContext uses ctx for HTTP requests.
func (ot *OpenTok) ForceDisconnectContext(ctx context.Context, sessionID, connectionID string) error {
	if sessionID == "" {
		return fmt.Errorf("Connection cannot be disconnected without a session ID")
	}

	if connectionID == "" {
		return fmt.Errorf("Connection cannot be disconnected without a connection ID")
	}

	// Create jwt token
	jwt, err := ot.genProjectJWT()
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/connection/" + connectionID
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

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
