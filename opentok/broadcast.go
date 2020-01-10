package opentok

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type HLSConfig struct{}

type RTMPConfig struct {
	Id         string `json:"id"`
	Status     string `json:"status,omitempty"`
	ServerUrl  string `json:"serverUrl"`
	StreamName string `json:"streamName"`
}

type BroadcastOutputOptions struct {
	HLS  *HLSConfig    `json:"hls,omitempty"`
	RTMP []*RTMPConfig `json:"rtmp,omitempty"`
}

type BroadcastOptions struct {
	SessionId   string                 `json:"sessionId"`
	Layout      *Layout                `json:"layout,omitempty"`
	MaxDuration int                    `json:"maxDuration,omitempty"`
	Outputs     BroadcastOutputOptions `json:"outputs"`
	Resolution  Resolution             `json:"resolution,omitempty"`
}

type BroadcastUrls struct {
	HLS  string        `json:"hls"`
	RTMP []*RTMPConfig `json:rtmp`
}

type Broadcast struct {
	Id            string        `json:"id"`            // The unique ID for the broadcast.
	SessionId     string        `json:"sessionId"`     // The OpenTok session ID.
	ProjectId     int           `json:"projectId"`     // The API key associated with the broadcast..
	CreatedAt     int           `json:"createdAt"`     // The time at which the broadcast was created, in milliseconds since the UNIX epoch.
	UpdatedAt     int           `json:"updatedAt"`     // The time at which the broadcast was updated, in milliseconds since the UNIX epoch.
	Resolution    Resolution    `json:"resolution"`    // The resolution of the broadcast.
	Status        string        `json:"status"`        // The status of the broadcast.
	BroadcastUrls BroadcastUrls `json:"broadcastUrls"` // An object containing details about the HLS and RTMP broadcasts.
	OpenTok       *OpenTok      `json:"-"`
}

type BroadcastListOptions struct {
	Offset    int
	Count     int
	SessionId string
}

type BroadcastList struct {
	Count int          `json:"count"`
	Items []*Broadcast `json:"items"`
}

/**
 * Use this method to start a live streaming for an OpenTok session.
 * This broadcasts the session to an HLS (HTTP live streaming) or to RTMP streams.
 */
func (ot *OpenTok) StartBroadcast(sessionId string, opts BroadcastOptions) (*Broadcast, error) {
	opts.SessionId = sessionId

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

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/broadcast"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
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

	broadcast := &Broadcast{}
	if err := json.NewDecoder(res.Body).Decode(broadcast); err != nil {
		return nil, err
	}

	broadcast.OpenTok = ot

	return broadcast, nil
}

/**
 * Use this method to stop a live broadcast of an OpenTok session.
 * Note that broadcasts automatically stop 120 minutes after they are started.
 */
func (ot *OpenTok) StopBroadcast(broadcastId string) (*Broadcast, error) {
	if broadcastId == "" {
		return nil, fmt.Errorf("Live stremaing broadcast cannot be stopped without an broadcast ID")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/broadcast/" + broadcastId + "/stop"
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

	broadcast := &Broadcast{}
	if err = json.NewDecoder(res.Body).Decode(broadcast); err != nil {
		return nil, err
	}

	broadcast.OpenTok = ot

	return broadcast, nil
}

/**
 * Get details on broadcasts that are in progress and started.
 * Completed broadcasts are not included in the listing.
 */
func (ot *OpenTok) ListBroadcasts(opts BroadcastListOptions) (*BroadcastList, error) {
	params := []string{"?"}

	if opts.Offset != 0 {
		params = append(params, "offset="+strconv.Itoa(opts.Offset))
	}

	if opts.Count != 0 {
		params = append(params, "count="+strconv.Itoa(opts.Count))
	}

	if opts.SessionId != "" {
		params = append(params, "sessionId="+opts.SessionId)
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/broadcast" + strings.Join(params, "&")
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

	broadcastList := &BroadcastList{}
	if err = json.NewDecoder(res.Body).Decode(&broadcastList); err != nil {
		return nil, err
	}

	for _, broadcast := range broadcastList.Items {
		broadcast.OpenTok = ot
	}

	return broadcastList, nil
}

/**
 * Get details on a broadcast that is in-progress.
 */
func (ot *OpenTok) GetBroadcast(broadcastId string) (*Broadcast, error) {
	if broadcastId == "" {
		return nil, fmt.Errorf("Cannot get broadcast information without an broadcast ID")
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/broadcast/" + broadcastId
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

	broadcast := &Broadcast{}
	if err := json.NewDecoder(res.Body).Decode(broadcast); err != nil {
		return nil, err
	}

	broadcast.OpenTok = ot

	return broadcast, nil
}

/**
 * Dynamically change the layout type of a live streaming broadcast.
 */
func (ot *OpenTok) SetBroadcastLayout(broadcastId string, layout Layout) (*Broadcast, error) {
	if broadcastId == "" {
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

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/broadcast/" + broadcastId + "/layout"
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

	broadcast := &Broadcast{}
	if err := json.NewDecoder(res.Body).Decode(broadcast); err != nil {
		return nil, err
	}

	broadcast.OpenTok = ot

	return broadcast, nil
}

func (broadcast *Broadcast) Stop() (*Broadcast, error) {
	return broadcast.OpenTok.StopBroadcast(broadcast.Id)
}
