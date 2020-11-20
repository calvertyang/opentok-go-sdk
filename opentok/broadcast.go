package opentok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// HLSConfig defines the config of HLS.
type HLSConfig struct{}

// RTMPConfig defines the config of RTMP.
type RTMPConfig struct {
	// An unique ID for the stream.
	ID string `json:"id"`

	// RTMP stream status returned in the response.
	Status string `json:"status,omitempty"`

	// The RTMP server URL.
	ServerURL string `json:"serverUrl"`

	// The stream name, such as the YouTube Live stream name or the Facebook
	// stream key.
	StreamName string `json:"streamName"`
}

// BroadcastOutputOptions defines the types of broadcast streams you want to
// start (both HLS and RTMP).
type BroadcastOutputOptions struct {
	// Set this property to an empty object for HLS.
	// The HLS URL is returned in the response and in the REST method for getting
	// information about a live streaming broadcast.
	HLS *HLSConfig `json:"hls,omitempty"`

	// The configuration of RTMP
	RTMP []*RTMPConfig `json:"rtmp,omitempty"`
}

// BroadcastOptions defines the options for live streaming broadcast.
type BroadcastOptions struct {
	// The OpenTok session you want to broadcast.
	SessionID string `json:"sessionId"`

	// Specify this to assign the initial layout type for the broadcast.
	Layout *Layout `json:"layout,omitempty"`

	// The maximum duration for the broadcast, in seconds.
	MaxDuration int `json:"maxDuration,omitempty"`

	// This object defines the types of broadcast streams you want to start
	// (both HLS and RTMP).
	Outputs *BroadcastOutputOptions `json:"outputs"`

	// The resolution of the broadcast: either SD(default) or HD.
	Resolution Resolution `json:"resolution,omitempty"`
}

// BroadcastURLs defines the details on the HLS and RTMP broadcast streams.
type BroadcastURLs struct {
	// HLS broadcast streams URL.
	HLS string `json:"hls"`

	// The configuration of RTMP.
	RTMP []*RTMPConfig `json:"rtmp"`
}

// Broadcast defines the response returned from API.
type Broadcast struct {
	// The unique ID for the broadcast.
	ID string `json:"id"`

	// The OpenTok session ID.
	SessionID string `json:"sessionId"`

	// The API key associated with the broadcast..
	ProjectID int `json:"projectId"`

	// The time at which the broadcast was created, in milliseconds since the
	// UNIX epoch.
	CreatedAt int `json:"createdAt"`

	// The time at which the broadcast was updated, in milliseconds since the
	// UNIX epoch.
	UpdatedAt int `json:"updatedAt"`

	// The resolution of the broadcast.
	Resolution Resolution `json:"resolution"`

	// The status of the broadcast.
	Status string `json:"status"`

	// An object containing details about the HLS and RTMP broadcasts.
	BroadcastURLs *BroadcastURLs `json:"broadcastUrls"`

	// The instance of OpenTok
	OpenTok *OpenTok `json:"-"`
}

// BroadcastListOptions defines the query parameters to filter the list of
// broadcasts.
type BroadcastListOptions struct {
	// The start offset in the list of existing broadcasts.
	Offset int

	// The number of broadcasts to retrieve starting at offset.
	Count int

	// Retrive only broadcasts for a given session ID.
	SessionID string
}

// BroadcastList defines the response returned from API.
type BroadcastList struct {
	// The total number of broadcasts in the results.
	Count int `json:"count"`

	// An array of objects defining each broadcast retrieved.
	// Broadcasts are listed from the newest to the oldest in the return set.
	Items []*Broadcast `json:"items"`
}

// StartBroadcast starts a live streaming for an OpenTok session. This
// broadcasts the session to an HLS (HTTP live streaming) or to RTMP
// streams.
func (ot *OpenTok) StartBroadcast(sessionID string, opts *BroadcastOptions) (*Broadcast, error) {
	return ot.StartBroadcastContext(context.Background(), sessionID, opts)
}

// StartBroadcastContext uses ctx for HTTP requests.
func (ot *OpenTok) StartBroadcastContext(ctx context.Context, sessionID string, opts *BroadcastOptions) (*Broadcast, error) {
	opts.SessionID = sessionID

	if opts.Layout != nil {
		if opts.Layout.Type != BestFit && opts.Layout.Type != PIP && opts.Layout.Type != Custom &&
			opts.Layout.Type != VerticalPresentation && opts.Layout.Type != HorizontalPresentation {
			return nil, fmt.Errorf("Invalid type of layout for starting a live streaming broadcast")
		}

		if opts.Layout.Type == Custom && opts.Layout.StyleSheet == "" {
			return nil, fmt.Errorf("StyleSheet property of layout cannot be empty")
		}

		// For other layout types, do not set a stylesheet property.
		if opts.Layout.Type != Custom && opts.Layout.StyleSheet != "" {
			return nil, fmt.Errorf("Set stylesheet property only when using custom layout")
		}
	}

	if opts.Resolution != "" && opts.Resolution != SD && opts.Resolution != HD {
		return nil, fmt.Errorf("Invalid resolution for starting a live streaming broadcast")
	}

	jsonStr, _ := json.Marshal(opts)

	// Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/broadcast"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
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

	broadcast := &Broadcast{}
	if err := json.NewDecoder(res.Body).Decode(broadcast); err != nil {
		return nil, err
	}

	broadcast.OpenTok = ot

	return broadcast, nil
}

// StopBroadcast stops a live broadcast of an OpenTok session.
// Note that broadcasts automatically stop 120 minutes after they are started.
func (ot *OpenTok) StopBroadcast(broadcastID string) (*Broadcast, error) {
	return ot.StopBroadcastContext(context.Background(), broadcastID)
}

// StopBroadcastContext uses ctx for HTTP requests.
func (ot *OpenTok) StopBroadcastContext(ctx context.Context, broadcastID string) (*Broadcast, error) {
	if broadcastID == "" {
		return nil, fmt.Errorf("Live stremaing broadcast cannot be stopped without an broadcast ID")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/broadcast/" + broadcastID + "/stop"
	req, err := http.NewRequest("POST", endpoint, nil)
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

	broadcast := &Broadcast{}
	if err = json.NewDecoder(res.Body).Decode(broadcast); err != nil {
		return nil, err
	}

	broadcast.OpenTok = ot

	return broadcast, nil
}

// ListBroadcasts returns the records of all broadcasts for your
// project that are in progress and started. Completed broadcasts are
// not included in the listing.
func (ot *OpenTok) ListBroadcasts(opts *BroadcastListOptions) (*BroadcastList, error) {
	return ot.ListBroadcastsContext(context.Background(), opts)
}

// ListBroadcastsContext uses ctx for HTTP requests.
func (ot *OpenTok) ListBroadcastsContext(ctx context.Context, opts *BroadcastListOptions) (*BroadcastList, error) {
	params := []string{"?"}

	if opts.Offset != 0 {
		params = append(params, "offset="+strconv.Itoa(opts.Offset))
	}

	if opts.Count != 0 {
		params = append(params, "count="+strconv.Itoa(opts.Count))
	}

	if opts.SessionID != "" {
		params = append(params, "sessionId="+opts.SessionID)
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/broadcast" + strings.Join(params, "&")
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

	broadcastList := &BroadcastList{}
	if err = json.NewDecoder(res.Body).Decode(&broadcastList); err != nil {
		return nil, err
	}

	for _, broadcast := range broadcastList.Items {
		broadcast.OpenTok = ot
	}

	return broadcastList, nil
}

// GetBroadcast returns a broadcast that is in-progress.
func (ot *OpenTok) GetBroadcast(broadcastID string) (*Broadcast, error) {
	return ot.GetBroadcastContext(context.Background(), broadcastID)
}

// GetBroadcastContext uses ctx for HTTP requests.
func (ot *OpenTok) GetBroadcastContext(ctx context.Context, broadcastID string) (*Broadcast, error) {
	if broadcastID == "" {
		return nil, fmt.Errorf("Cannot get broadcast information without an broadcast ID")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/broadcast/" + broadcastID
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

	broadcast := &Broadcast{}
	if err := json.NewDecoder(res.Body).Decode(broadcast); err != nil {
		return nil, err
	}

	broadcast.OpenTok = ot

	return broadcast, nil
}

// SetBroadcastLayout dynamically changes the layout type of a live streaming
// broadcast.
func (ot *OpenTok) SetBroadcastLayout(broadcastID string, layout *Layout) (*Broadcast, error) {
	return ot.SetBroadcastLayoutContext(context.Background(), broadcastID, layout)
}

// SetBroadcastLayoutContext uses ctx for HTTP requests.
func (ot *OpenTok) SetBroadcastLayoutContext(ctx context.Context, broadcastID string, layout *Layout) (*Broadcast, error) {
	if broadcastID == "" {
		return nil, fmt.Errorf("Cannot change the layout type of a live streaming broadcast without an broadcast ID")
	}

	if layout.Type != BestFit && layout.Type != PIP && layout.Type != Custom &&
		layout.Type != VerticalPresentation && layout.Type != HorizontalPresentation {
		return nil, fmt.Errorf("Invalid type of layout for archive")
	}

	if layout.Type == Custom && layout.StyleSheet == "" {
		return nil, fmt.Errorf("StyleSheet property of layout cannot be empty")
	}

	// For other layout types, do not set a stylesheet property.
	if layout.Type != Custom && layout.StyleSheet != "" {
		return nil, fmt.Errorf("Set stylesheet property only when using custom layout")
	}

	jsonStr, _ := json.Marshal(layout)

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/broadcast/" + broadcastID + "/layout"
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

	broadcast := &Broadcast{}
	if err := json.NewDecoder(res.Body).Decode(broadcast); err != nil {
		return nil, err
	}

	broadcast.OpenTok = ot

	return broadcast, nil
}

// Stop stops a live broadcast of an OpenTok session.
func (broadcast *Broadcast) Stop() (*Broadcast, error) {
	return broadcast.OpenTok.StopBroadcast(broadcast.ID)
}
