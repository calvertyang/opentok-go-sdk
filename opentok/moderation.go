package opentok

import (
	"fmt"
	"net/http"
)

/**
 * Disconnect a client from an OpenTok session via server-side
 */
func (ot *OpenTok) ForceDisconnect(sessionId, connectionId string) error {
	if sessionId == "" {
		return fmt.Errorf("Connection cannot be disconnected without a session ID")
	}

	if connectionId == "" {
		return fmt.Errorf("Connection cannot be disconnected without a connection ID")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return err
	}

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionId + "/connection/" + connectionId
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
