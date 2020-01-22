package opentok

import (
	"fmt"
	"net/http"
)

// ForceDisconnect disconnects a client from an OpenTok session via server-side
func (ot *OpenTok) ForceDisconnect(sessionID, connectionID string) error {
	if sessionID == "" {
		return fmt.Errorf("Connection cannot be disconnected without a session ID")
	}

	if connectionID == "" {
		return fmt.Errorf("Connection cannot be disconnected without a connection ID")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/connection/" + connectionID
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

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
