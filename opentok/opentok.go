package opentok

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// HTTPClient is an interface to allow custom clients and timeouts.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// OpenTok API host URL
const defaultAPIHost = "https://api.opentok.com"

// For use in X-TB-TOKEN-AUTH header value
const tokenSentinel = "T1=="

type issueType string

const (
	// For most REST API calls, set issue type to "project"
	projectToken issueType = "project"
	// For Account Management REST methods, set issue type to "account"
	accountToken issueType = "account"
)

// OpenTok stores the API key and secret for use in making API call
type OpenTok struct {
	apiKey    string
	apiSecret string
	apiHost   string

	httpClient HTTPClient
}

// New returns an initialized OpenTok instance with the API key and API secret.
func New(apiKey, apiSecret string, client HTTPClient) *OpenTok {
	return &OpenTok{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		apiHost:    defaultAPIHost,
		httpClient: client,
	}
}

// SetAPIHost is used to set OpenTok API Host to specific URL
func (ot *OpenTok) SetAPIHost(url string) error {
	if url == "" {
		return fmt.Errorf("OpenTok API Host cannot be empty")
	}

	ot.apiHost = url

	return nil
}

// Generate JWT token for API calls
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
