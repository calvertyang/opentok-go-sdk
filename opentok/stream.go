package opentok

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Stream struct {
	Id              string   `json:"id"`              // The unique stream ID.
	VideoType       string   `json:"videoType"`       // Either "camera" or "screen".
	Name            string   `json:"name"`            // The stream name.
	LayoutClassList []string `json:"layoutClassList"` // An array of the layout classes for the stream.
}

type StreamList struct {
	Count int       `json:"count"`
	Items []*Stream `json:"items"`
}

/**
 * Get information on an OpenTok all stream in a session
 */
func (ot *OpenTok) ListStreams(sessionId string) (*StreamList, error) {
	if sessionId == "" {
		return nil, fmt.Errorf("Cannot get all streams information without a session ID")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionId + "/stream"
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

	streamList := &StreamList{}
	if err := json.NewDecoder(res.Body).Decode(streamList); err != nil {
		return nil, err
	}

	return streamList, nil
}

/**
 * Get information on an OpenTok stream in a session
 */
func (ot *OpenTok) GetStream(sessionId, streamId string) (*Stream, error) {
	if sessionId == "" {
		return nil, fmt.Errorf("Cannot get stream information without a session ID")
	}

	if streamId == "" {
		return nil, fmt.Errorf("Cannot get stream information without a stream ID")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionId + "/stream/" + streamId
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

	stream := &Stream{}
	if err := json.NewDecoder(res.Body).Decode(stream); err != nil {
		return nil, err
	}

	return stream, nil
}
