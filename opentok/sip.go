package opentok

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SIPHeaders map[string]string

type SIPAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SIP struct {
	URI     string      `json:"uri"`
	From    string      `json:"from,omitempty"`
	Headers *SIPHeaders `json:"headers,omitempty"`
	Auth    *SIPAuth    `json:"auth,omitempty"`
	Secure  bool        `json:"secure,omitempty"`
}

type DialOptions struct {
	SessionId string `json:"sessionId"`
	Token     string `json:"token"`
	SIP       SIP    `json:"sip"`
	TokenData string `json:"-"`
}

type SIPCall struct {
	Id           string `json:"id"`           // A unique ID for the SIP call.
	ConnectionId string `json:"connectionId"` // The OpenTok connection ID for the SIP call's connection in the OpenTok session.
	StreamId     string `json:"streamId"`     // The OpenTok stream ID for the SIP call's stream in the OpenTok session.
}

/**
 * The audio from your end of the SIP call is added to the OpenTok session as
 * an audio-only stream. The OpenTok Media Router mixes audio from other streams
 * in the session and sends the mixed audio to your SIP endpoint.
 */
func (ot *OpenTok) Dial(sessionId string, opts DialOptions) (*SIPCall, error) {
	if sessionId == "" {
		return nil, fmt.Errorf("SIP call cannot be initiated without an session ID")
	}

	if opts.SIP.URI == "" {
		return nil, fmt.Errorf("SIP call cannot be initiated without an SIP URI")
	}

	token, err := ot.GenerateToken(sessionId, TokenOptions{
		Data: opts.TokenData,
	})
	if err != nil {
		return nil, err
	}

	opts.SessionId = sessionId
	opts.Token = *token

	jsonStr, _ := json.Marshal(opts)

	//Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := apiHost + projectURL + "/" + ot.apiKey + "/dial"
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

	sipCall := &SIPCall{}
	if err := json.NewDecoder(res.Body).Decode(sipCall); err != nil {
		return nil, err
	}

	return sipCall, nil
}
