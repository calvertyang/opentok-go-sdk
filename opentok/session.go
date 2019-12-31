package opentok

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const sessionCreateURL = "/session/create"

type ArchiveMode string

const (
	/**
	 * Set to always to have the session archived automatically.
	 */
	AutoArchived ArchiveMode = "always"
	/**
	 * Set to manual (the default), you can archive the session by calling the REST /archive POST method
	 */
	ManualArchived ArchiveMode = "manual"
)

type MediaMode string

const (
	/**
	 * Set to enabled if you prefer clients to attempt to send audio-video streams directly to other clients
	 */
	Relayed MediaMode = "enabled"
	/**
	 * Set to disabled for sessions that use the OpenTok Media Router
	 */
	Routed MediaMode = "disabled"
)

type SessionOptions struct {
	ArchiveMode ArchiveMode
	Location    string
	MediaMode   MediaMode
}

type Session struct {
	SessionId      string   `json:"session_id"`
	ProjectId      string   `json:"project_id"`
	CreateDt       string   `json:"create_dt"`
	MediaServerURL string   `json:"media_server_url"`
	OpenTok        *OpenTok `json:"-"`
}

func (ot *OpenTok) CreateSession(opts SessionOptions) (*Session, error) {
	params := url.Values{}

	if opts.ArchiveMode != "" {
		params.Add("archiveMode", string(opts.ArchiveMode))
	}

	if opts.Location != "" {
		params.Add("location", opts.Location)
	}

	if opts.MediaMode != "" {
		params.Add("p2p.preference", string(opts.MediaMode))
	}

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiHost+sessionCreateURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
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

	var sessions []Session
	if err = json.NewDecoder(res.Body).Decode(&sessions); err != nil {
		return nil, err
	}

	if len(sessions) < 1 {
		return nil, fmt.Errorf("Tokbox did not return a session")
	}

	session := sessions[0]
	session.OpenTok = ot

	return &session, nil
}

func (s *Session) GenerateToken(opts TokenOptions) (*string, error) {
	return s.OpenTok.GenerateToken(s.SessionId, opts)
}
