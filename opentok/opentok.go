package opentok

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const apiHost = "https://api.opentok.com"

const TOKEN_SENTINEL = "T1=="

type issueType string

const (
	/**
	 * For most REST API calls, set issue type to "project"
	 */
	projectToken issueType = "project"
	/**
	 * For Account Management REST methods, set issue type to "account"
	 */
	accountToken issueType = "account"
)

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

type OpenTok struct {
	apiKey    string
	apiSecret string
}

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

func New(apiKey, apiSecret string) *OpenTok {
	return &OpenTok{apiKey, apiSecret}
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

func (ot *OpenTok) jwtToken(ist issueType) (string, error) {
	type OpenTokClaims struct {
		Ist issueType `json:"ist,omitempty"`
		jwt.StandardClaims
	}

	issuedAt := time.Now().UTC()

	claims := OpenTokClaims{
		ist,
		jwt.StandardClaims{
			Issuer:    ot.apiKey,
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: issuedAt.Add((5 * time.Minute)).Unix(), // The maximum allowed expiration time range is 5 minutes.
			Id:        uuid.New().String(),
		},
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the api secret
	return token.SignedString([]byte(ot.apiSecret))
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
