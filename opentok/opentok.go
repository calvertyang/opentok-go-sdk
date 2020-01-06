package opentok

import (
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

type OpenTok struct {
	apiKey    string
	apiSecret string
}

func New(apiKey, apiSecret string) *OpenTok {
	return &OpenTok{apiKey, apiSecret}
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
