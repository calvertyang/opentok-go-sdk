package opentok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Stream defines the response returned from API
type Stream struct {
	// The unique stream ID.
	ID string `json:"id"`

	// Either "camera" or "screen".
	VideoType string `json:"videoType"`

	// The stream name.
	Name string `json:"name"`

	// An array of the layout classes for the stream.
	LayoutClassList []string `json:"layoutClassList"`
}

// StreamList defines the response returned from API
type StreamList struct {
	// The total number of streams in a session.
	Count int `json:"count"`

	// An array of objects defining each stream retrieved.
	// Streams are listed from the newest to the oldest in the return set.
	Items []*Stream `json:"items"`
}

// StreamClass defines the layout classes to assign to a stream.
type StreamClass struct {
	// The stream ID.
	ID string `json:"id"`

	// An array of layout classes (each strings) for the stream.
	LayoutClassList []string `json:"layoutClassList"`
}

// StreamClassOptions defines the options for setting the layout classes for
// the stream
type StreamClassOptions struct {
	// The layout classes to assign to a stream.
	Items []*StreamClass `json:"items"`
}

// ListStreams returns the stream records in a session.
func (ot *OpenTok) ListStreams(sessionID string) (*StreamList, error) {
	return ot.ListStreamsContext(context.Background(), sessionID)
}

// ListStreamsContext uses ctx for HTTP requests.
func (ot *OpenTok) ListStreamsContext(ctx context.Context, sessionID string) (*StreamList, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("Cannot get all streams information without a session ID")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/stream"
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

	streamList := &StreamList{}
	if err := json.NewDecoder(res.Body).Decode(streamList); err != nil {
		return nil, err
	}

	return streamList, nil
}

// GetStream returns a stream details record describing the stream.
func (ot *OpenTok) GetStream(sessionID, streamID string) (*Stream, error) {
	return ot.GetStreamContext(context.Background(), sessionID, streamID)
}

// GetStreamContext uses ctx for HTTP requests.
func (ot *OpenTok) GetStreamContext(ctx context.Context, sessionID, streamID string) (*Stream, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("Cannot get stream information without a session ID")
	}

	if streamID == "" {
		return nil, fmt.Errorf("Cannot get stream information without a stream ID")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/stream/" + streamID
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

	stream := &Stream{}
	if err := json.NewDecoder(res.Body).Decode(stream); err != nil {
		return nil, err
	}

	return stream, nil
}

// SetStreamClassLists changes the composed archive layout classes for an OpenTok stream
func (ot *OpenTok) SetStreamClassLists(sessionID string, opts *StreamClassOptions) (*StreamList, error) {
	return ot.SetStreamClassListsContext(context.Background(), sessionID, opts)
}

// SetStreamClassListsContext uses ctx for HTTP requests.
func (ot *OpenTok) SetStreamClassListsContext(ctx context.Context, sessionID string, opts *StreamClassOptions) (*StreamList, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("Cannot change the live streaming layout classes for an OpenTok stream without an session ID")
	}

	jsonStr, _ := json.Marshal(opts)

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/stream"
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

	streamList := &StreamList{}
	if err := json.NewDecoder(res.Body).Decode(streamList); err != nil {
		return nil, err
	}

	return streamList, nil
}
