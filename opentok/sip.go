package opentok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// SIPHeaders is the alias of map[string]string type
type SIPHeaders map[string]string

// SIPAuth defines the authentication information for SIP call
type SIPAuth struct {
	// The username for HTTP digest authentication
	Username string `json:"username"`

	// The password for HTTP digest authentication
	Password string `json:"password"`
}

// SIP defines the information of SIP call
type SIP struct {
	// The SIP URI to be used as destination of the SIP call initiated from
	// OpenTok to your SIP platform.
	URI string `json:"uri"`

	// The number or string that will be sent to the final SIP number as the
	// caller.
	From string `json:"from,omitempty"`

	// Custom headers to be added to the SIP ​INVITE​ request initiated from
	// OpenTok to your SIP platform.
	Headers *SIPHeaders `json:"headers,omitempty"`

	// The username and password to be used in the the SIP INVITE​ request for
	// HTTP digest authentication
	Auth *SIPAuth `json:"auth,omitempty"`

	// A flag that indicates whether the media must be transmitted encrypted
	// or not.
	Secure bool `json:"secure,omitempty"`
}

// DialOptions defines the options for SIP call
type DialOptions struct {
	// The OpenTok session ID for the SIP call to join.
	SessionID string `json:"sessionId"`

	// The OpenTok token to be used for the participant being called.
	Token string `json:"token"`

	// The SIP information
	SIP *SIP `json:"sip"`

	// The data for token generation
	TokenData string `json:"-"`
}

// SIPCall defines the response returned from API
type SIPCall struct {
	// A unique ID for the SIP call.
	ID string `json:"id"`

	// The OpenTok connection ID for the SIP call's connection in the OpenTok
	// session.
	ConnectionID string `json:"connectionId"`

	// The OpenTok stream ID for the SIP call's stream in the OpenTok session.
	StreamID string `json:"streamId"`
}

// Dial connects your SIP platform to an OpenTok session.
// The audio from your end of the SIP call is added to the OpenTok session as
// an audio-only stream. The OpenTok Media Router mixes audio from other streams
// in the session and sends the mixed audio to your SIP endpoint.
func (ot *OpenTok) Dial(sessionID string, opts *DialOptions) (*SIPCall, error) {
	return ot.DialContext(context.Background(), sessionID, opts)
}

// DialContext uses ctx for HTTP requests.
func (ot *OpenTok) DialContext(ctx context.Context, sessionID string, opts *DialOptions) (*SIPCall, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("SIP call cannot be initiated without a session ID")
	}

	if opts.SIP.URI == "" {
		return nil, fmt.Errorf("SIP call cannot be initiated without a SIP URI")
	}

	token, err := ot.GenerateToken(sessionID, &TokenOptions{
		Data: opts.TokenData,
	})
	if err != nil {
		return nil, err
	}

	opts.SessionID = sessionID
	opts.Token = token

	jsonStr, _ := json.Marshal(opts)

	// Create jwt token
	jwt, err := ot.genProjectJWT()
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/dial"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)

	res, err := ot.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, parseErrorResponse(res)
	}

	sipCall := &SIPCall{}
	if err := json.NewDecoder(res.Body).Decode(sipCall); err != nil {
		return nil, err
	}

	return sipCall, nil
}

// SendDTMF sends the DTMF digits to all clients connected to the session.
func (ot *OpenTok) SendDTMF(sessionID string, digits string) error {
	return ot.SendDTMFContext(context.Background(), sessionID, digits)
}

// SendDTMFContext uses ctx for HTTP requests.
func (ot *OpenTok) SendDTMFContext(ctx context.Context, sessionID, digits string) error {
	if sessionID == "" {
		return fmt.Errorf("DTMF digits cannot be sent without a session ID")
	}

	if digits == "" {
		return fmt.Errorf("The DTMF digits cannot be empty")
	}

	match, err := regexp.MatchString(`^[\d#*p]+$`, digits)
	if !match || err != nil {
		return fmt.Errorf("The DTMF digits is invalid")
	}

	jsonStr := []byte(`{ "digits": "` + digits + `" }`)

	// Create jwt token
	jwt, err := ot.genProjectJWT()
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/play-dtmf"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)

	res, err := ot.sendRequest(ctx, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return parseErrorResponse(res)
	}

	return nil
}

// SendDTMFToClient sends the DTMF tones to a specific client connected to the session.
func (ot *OpenTok) SendDTMFToClient(sessionID, connectionID, digits string) error {
	return ot.SendDTMFToClientContext(context.Background(), sessionID, connectionID, digits)
}

// SendDTMFToClientContext uses ctx for HTTP requests.
func (ot *OpenTok) SendDTMFToClientContext(ctx context.Context, sessionID, connectionID, digits string) error {
	if sessionID == "" {
		return fmt.Errorf("DTMF digits cannot be sent without a session ID")
	}

	if connectionID == "" {
		return fmt.Errorf("DTMF digits cannot be sent without a connection ID")
	}

	if digits == "" {
		return fmt.Errorf("The DTMF digits cannot be empty")
	}

	match, err := regexp.MatchString(`^[\d#*p]+$`, digits)
	if !match || err != nil {
		return fmt.Errorf("The DTMF digits is invalid")
	}

	jsonStr := []byte(`{ "digits": "` + digits + `" }`)

	// Create jwt token
	jwt, err := ot.genProjectJWT()
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/session/" + sessionID + "/connection/" + connectionID + "/play-dtmf"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)

	res, err := ot.sendRequest(ctx, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return parseErrorResponse(res)
	}

	return nil
}

// Dial connects your SIP platform to an OpenTok session.
func (s *Session) Dial(opts *DialOptions) (*SIPCall, error) {
	return s.DialContext(context.Background(), opts)
}

// DialContext uses ctx for HTTP requests.
func (s *Session) DialContext(ctx context.Context, opts *DialOptions) (*SIPCall, error) {
	return s.OpenTok.DialContext(ctx, s.SessionID, opts)
}

// SendDTMF sends the DTMF digits to all clients connected to the session.
func (s *Session) SendDTMF(digits string) error {
	return s.SendDTMFContext(context.Background(), digits)
}

// SendDTMFContext uses ctx for HTTP requests.
func (s *Session) SendDTMFContext(ctx context.Context, digits string) error {
	return s.OpenTok.SendDTMFContext(ctx, s.SessionID, digits)
}

// SendDTMFToClient sends the DTMF tones to a specific client connected to the session.
func (s *Session) SendDTMFToClient(connectionID, digits string) error {
	return s.SendDTMFToClientContext(context.Background(), connectionID, digits)
}

// SendDTMFToClientContext uses ctx for HTTP requests.
func (s *Session) SendDTMFToClientContext(ctx context.Context, connectionID, digits string) error {
	return s.OpenTok.SendDTMFToClientContext(ctx, s.SessionID, connectionID, digits)
}
