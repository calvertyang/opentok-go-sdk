package opentok

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const sessionCreateURL = "/session/create"

// ArchiveMode is the alias of string type.
type ArchiveMode string

const (
	// AutoArchived is used to have the session archived automatically.
	AutoArchived ArchiveMode = "always"

	// ManualArchived is used to have the session archived manually, you can
	// archive the session by calling the REST /archive POST method.
	ManualArchived ArchiveMode = "manual"
)

// MediaMode is the alias of string type.
type MediaMode string

const (
	// Relayed if you prefer clients to attempt to send audio-video
	// streams directly to other clients.
	Relayed MediaMode = "enabled"

	// Routed if you want to use the OpenTok Media Router.
	Routed MediaMode = "disabled"
)

// SessionOptions defines the options for creating a session.
type SessionOptions struct {
	// Set to always to have the session archived automatically.
	// With the archiveModeset to manual (the default), you can archive the
	// session by calling the REST /archive POST method.
	ArchiveMode ArchiveMode

	// The IP address that TokBox will use to situate the session in its global
	// network.
	Location string

	// Set to enabled if you prefer clients to attempt to send audio-video
	// streams directly to other clients; set to disabled for sessions that use
	// the OpenTok Media Router.
	MediaMode MediaMode
}

// Session defines the response returned from API.
type Session struct {
	// The session ID.
	SessionID string `json:"session_id"`

	// The API key associated with the project.
	ProjectID string `json:"project_id"`

	// The time at which the session was created.
	CreateDt string `json:"create_dt"`

	// The URL of the OpenTok media router used by the session.
	MediaServerURL string `json:"media_server_url"`

	// The instance of OpenTok.
	OpenTok *OpenTok `json:"-"`
}

// Role is the alias of string type.
type Role string

const (
	// Publisher can publish streams, subscribe to streams, and signal.
	Publisher Role = "publisher"

	// Subscriber can only subscribe to streams.
	Subscriber Role = "subscriber"

	// Moderator can call the forceUnpublish() and forceDisconnect() method of
	// the Session object in clients using the OpenTok.js library, and have the
	// privileges granted to a publisher.
	Moderator Role = "moderator"
)

// TokenOptions defines the options for generating token.
type TokenOptions struct {
	// The role to determine the capabilities of the client that connects with
	// a token.
	Role Role

	// The metadata for describing the client.
	Data string

	// The expiration period of the token.
	ExpireTime int64

	// Layout classes for the stream.
	InitialLayoutClassList []string
}

// SessionIDInfo defines the information decoded from the session ID.
type SessionIDInfo struct {
	// The API key associated with the project.
	APIKey string

	// The IP address that TokBox use to situate the session in its global
	// network.
	Location string

	// The time at which the session was created.
	CreateTime time.Time
}

// CreateSession generates a new session.
func (ot *OpenTok) CreateSession(opts *SessionOptions) (*Session, error) {
	return ot.CreateSessionContext(context.Background(), opts)
}

// CreateSessionContext uses ctx for HTTP requests.
func (ot *OpenTok) CreateSessionContext(ctx context.Context, opts *SessionOptions) (*Session, error) {
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

	// Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ot.apiHost+sessionCreateURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
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

// GenerateToken generates a token for each user connecting to an OpenTok
// session.
func (ot *OpenTok) GenerateToken(sessionID string, opts *TokenOptions) (string, error) {
	if sessionID == "" {
		return "", fmt.Errorf("Token cannot be generated without a sessionID")
	}

	// validate the sessionID belongs to the apiKey of this OpenTok instance
	if sessionIDInfo, err := decodeSessionID(sessionID); err != nil || sessionIDInfo.APIKey != ot.apiKey {
		return "", fmt.Errorf("Token cannot be generated unless the session belongs to the API Key")
	}

	// create tokenData with given opts
	now := time.Now().UTC()
	rand.Seed(time.Now().UTC().UnixNano())
	tokenData := map[string]string{
		"session_id":                sessionID,
		"create_time":               strconv.FormatInt(now.Unix(), 10),
		"expire_time":               strconv.FormatInt(now.Add(24*time.Hour).Unix(), 10),
		"nonce":                     fmt.Sprintf("%v", rand.Float64()),
		"role":                      string(Publisher),
		"connection_data":           "",
		"initial_layout_class_list": "",
	}

	if opts.Role != "" {
		tokenData["role"] = string(opts.Role)
	}

	if opts.Data != "" {
		tokenData["connection_data"] = opts.Data
	}

	if opts.ExpireTime > 0 {
		tokenData["expire_time"] = strconv.FormatInt(opts.ExpireTime, 10)
	}

	if len(opts.InitialLayoutClassList) > 0 {
		tokenData["initial_layout_class_list"] = strings.Join(opts.InitialLayoutClassList, ",")
	}

	// validate tokenData
	if tokenData["role"] != string(Publisher) && tokenData["role"] != string(Subscriber) && tokenData["role"] != string(Moderator) {
		return "", fmt.Errorf("Invalid role for token generation: %v", tokenData["role"])
	}

	if tokenData["expire_time"] < tokenData["create_time"] {
		return "", fmt.Errorf("Invalid expireTime for token generation, time cannot be in the past: %v < %v", tokenData["expire_time"], tokenData["create_time"])
	}

	if tokenData["connection_data"] != "" && len(tokenData["connection_data"]) > 1024 {
		return "", fmt.Errorf("Invalid data for token generation, must be a string with maximum length 1024")
	}

	if tokenData["initial_layout_class_list"] != "" && len(tokenData["initial_layout_class_list"]) > 1024 {
		return "", fmt.Errorf("Invalid initial layout class list for token generation, must have concatenated length of less than 1024")
	}

	return encodeToken(tokenData, ot)
}

// GenerateToken generates a token for each user connecting to an OpenTok
// session.
func (s *Session) GenerateToken(opts *TokenOptions) (string, error) {
	return s.OpenTok.GenerateToken(s.SessionID, opts)
}

// Decodes a sessionID into the metadata that it contains
func decodeSessionID(sessionID string) (*SessionIDInfo, error) {
	// remove sentinel
	sessionID = sessionID[2:]

	// replace invalid base64 chars
	sessionID = strings.ReplaceAll(sessionID, "-", "+")
	sessionID = strings.ReplaceAll(sessionID, "_", "/")

	// base64 decode
	decodedSessionID, _ := base64.StdEncoding.DecodeString(sessionID)

	// separate fields
	fields := strings.Split(string(decodedSessionID), "~")

	ts, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return nil, err
	}

	sec := ts / 1000
	nsec := ts % 1000 * 1000000

	sessionIDInfo := &SessionIDInfo{
		APIKey:     fields[1],
		Location:   fields[2],
		CreateTime: time.Unix(sec, nsec),
	}

	return sessionIDInfo, nil
}

// Encodes data for use as a token that can be used as the X-TB-TOKEN-AUTH header value in OpenTok REST APIs
func encodeToken(tokenData map[string]string, ot *OpenTok) (string, error) {
	params := url.Values{}
	params.Add("session_id", tokenData["session_id"])
	params.Add("create_time", tokenData["create_time"])
	params.Add("expire_time", tokenData["expire_time"])
	params.Add("nonce", tokenData["nonce"])
	params.Add("role", tokenData["role"])
	params.Add("connection_data", tokenData["connection_data"])
	params.Add("initial_layout_class_list", tokenData["initial_layout_class_list"])

	dataString := params.Encode()

	h := hmac.New(sha1.New, []byte(ot.apiSecret))
	n, err := h.Write([]byte(dataString))
	if err != nil {
		return "", err
	}

	if n != len(dataString) {
		return "", fmt.Errorf("hmac not enough bytes written %d != %d", n, len(dataString))
	}

	sig := fmt.Sprintf("%x:%s", h.Sum(nil), dataString)
	decoded := "partner_id=" + ot.apiKey + "&sig=" + sig
	token := tokenSentinel + base64.StdEncoding.EncodeToString([]byte(decoded))

	return token, nil
}
