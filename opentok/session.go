package opentok

import (
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

type Role string

const (
	/**
	 * A publisher can publish streams, subscribe to streams, and signal.
	 */
	Publisher Role = "publisher"
	/**
	 * A subscriber can only subscribe to streams.
	 */
	Subscriber Role = "subscriber"
	/**
	 * In addition to the privileges granted to a publisher, in clients using
	 * the OpenTok.js library, a moderator can call the forceUnpublish() and
	 * forceDisconnect() method of the Session object.
	 */
	Moderator Role = "moderator"
)

type TokenOptions struct {
	Role                   Role
	Data                   string
	ExpireTime             int
	InitialLayoutClassList []string
}

type SessionIdInfo struct {
	ApiKey     string
	Location   string
	CreateTime time.Time
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

func (ot *OpenTok) GenerateToken(sessionId string, opts TokenOptions) (*string, error) {
	if sessionId == "" {
		return nil, fmt.Errorf("Token cannot be generated without a sessionId")
	}

	// validate the sessionId belongs to the apiKey of this OpenTok instance
	if sessionIdInfo, err := decodeSessionId(sessionId); err != nil || sessionIdInfo.ApiKey != ot.apiKey {
		return nil, fmt.Errorf("Token cannot be generated unless the session belongs to the API Key")
	}

	// create tokenData with given opts
	now := time.Now().UTC()
	rand.Seed(time.Now().UTC().UnixNano())
	tokenData := map[string]string{
		"session_id":                sessionId,
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
		tokenData["expire_time"] = strconv.Itoa(opts.ExpireTime)
	}

	if len(opts.InitialLayoutClassList) > 0 {
		tokenData["initial_layout_class_list"] = strings.Join(opts.InitialLayoutClassList, ",")
	}

	// validate tokenData
	if tokenData["role"] != string(Publisher) && tokenData["role"] != string(Subscriber) && tokenData["role"] != string(Moderator) {
		return nil, fmt.Errorf("Invalid role for token generation: %v", tokenData["role"])
	}

	if tokenData["expire_time"] < tokenData["create_time"] {
		return nil, fmt.Errorf("Invalid expireTime for token generation, time cannot be in the past: %v < %v", tokenData["expire_time"], tokenData["create_time"])
	}

	if tokenData["connection_data"] != "" && len(tokenData["connection_data"]) > 1024 {
		return nil, fmt.Errorf("Invalid data for token generation, must be a string with maximum length 1024")
	}

	if tokenData["initial_layout_class_list"] != "" && len(tokenData["initial_layout_class_list"]) > 1024 {
		return nil, fmt.Errorf("Invalid initial layout class list for token generation, must have concatenated length of less than 1024")
	}

	return encodeToken(tokenData, ot)
}

func (s *Session) GenerateToken(opts TokenOptions) (*string, error) {
	return s.OpenTok.GenerateToken(s.SessionId, opts)
}

/*
 * decodes a sessionId into the metadata that it contains
 */
func decodeSessionId(sessionId string) (*SessionIdInfo, error) {
	// remove sentinel
	sessionId = sessionId[2:]

	// replace invalid base64 chars
	sessionId = strings.ReplaceAll(sessionId, "-", "+")
	sessionId = strings.ReplaceAll(sessionId, "_", "/")

	// base64 decode
	decodedSessionId, _ := base64.StdEncoding.DecodeString(sessionId)

	// separate fields
	fields := strings.Split(string(decodedSessionId), "~")

	ts, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return nil, err
	}

	sec := ts / 1000
	nsec := ts % 1000 * 1000000

	sessionIdInfo := &SessionIdInfo{
		ApiKey:     fields[1],
		Location:   fields[2],
		CreateTime: time.Unix(sec, nsec),
	}

	return sessionIdInfo, nil
}

/**
 * Encodes data for use as a token that can be used as the X-TB-TOKEN-AUTH header value in OpenTok REST APIs
 */
func encodeToken(tokenData map[string]string, ot *OpenTok) (*string, error) {
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
		return nil, err
	}

	if n != len(dataString) {
		return nil, fmt.Errorf("hmac not enough bytes written %d != %d", n, len(dataString))
	}

	sig := fmt.Sprintf("%x:%s", h.Sum(nil), dataString)
	decoded := "partner_id=" + ot.apiKey + "&sig=" + sig
	token := TOKEN_SENTINEL + base64.StdEncoding.EncodeToString([]byte(decoded))

	return &token, nil
}
