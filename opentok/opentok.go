package opentok

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// OpenTok API host URL.
const defaultAPIHost = "https://api.opentok.com"

// For use in X-TB-TOKEN-AUTH header value.
const tokenSentinel = "T1=="

type issueType string

const (
	// For most REST API calls, set issue type to "project".
	projectToken issueType = "project"

	// For Account Management REST methods, set issue type to "account".
	accountToken issueType = "account"
)

// HTTPClient is an interface to allow custom clients and timeouts.
type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

// OpenTokResponseError encloses an error with code and message.
type OpenTokResponseError struct {
	// StatusCode is the HTTP Response StatusCode that led to the error.
	StatusCode int

	// Message is the error message.
	Message string `json:"message"`
}

// Error returns a formatted error message.
func (e *OpenTokResponseError) Error() string {
	return fmt.Sprintf("TokBox error: code: %d; message: %s", e.StatusCode, e.Message)
}

// OpenTok stores the API key and secret for use in making API call.
type OpenTok struct {
	apiKey    string
	apiSecret string
	apiHost   string
	debug     bool

	httpClient HTTPClient
}

// New returns an initialized OpenTok instance with the API key and API secret.
func New(apiKey, apiSecret string) *OpenTok {
	return &OpenTok{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		apiHost:    defaultAPIHost,
		debug:      false,
		httpClient: http.DefaultClient,
	}
}

// SetAPIHost is used to set the OpenTok API Host to specific URL.
func (ot *OpenTok) SetAPIHost(url string) error {
	if url == "" {
		return fmt.Errorf("OpenTok API Host cannot be empty")
	}

	ot.apiHost = url

	return nil
}

// Debug sets the debug flage to true.
// It will dump the request/response allows you easy to debug.
func (ot *OpenTok) Debug() {
	ot.debug = true
}

// SetHttpClient specifies http client, http.DefaultClient used by default.
func (ot *OpenTok) SetHttpClient(client HTTPClient) {
	if client != nil {
		ot.httpClient = client
	}
}

// Generate account-level JWT token for API calls.
func (ot *OpenTok) genAccountJWT() (string, error) {
	return ot.jwtToken(accountToken)
}

// Generate project-level JWT token for API calls.
func (ot *OpenTok) genProjectJWT() (string, error) {
	return ot.jwtToken(projectToken)
}

// Generate JWT token for API calls.
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

// Send HTTP request.
func (ot *OpenTok) sendRequest(req *http.Request, ctx context.Context) (*http.Response, error) {
	req.Header.Add("User-Agent", userAgent)

	// Dump request
	if ot.debug {
		fmt.Println("========== Request Begin ==========")
		reqDump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(reqDump))
		fmt.Println("========== Request End ==========")
	}

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	// Dump response
	if ot.debug {
		fmt.Println("========== Response Begin ==========")
		resDump, _ := httputil.DumpResponse(res, true)
		fmt.Println(string(resDump))
		fmt.Println("========== Response End ==========")
	}

	return res, err
}

// Parse the error rresponse by custom error struct.
func parseErrorResponse(res *http.Response) error {
	resErr := &OpenTokResponseError{}
	if err := json.NewDecoder(res.Body).Decode(resErr); err != nil {
		return fmt.Errorf("Error decoding response from Tokbox: statusCode: %d; %w", res.StatusCode, err)
	}

	resErr.StatusCode = res.StatusCode

	return resErr
}
