package opentok

import (
	"log"
	"net/http"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

const (
	apiKey    = "<your api key here>"
	apiSecret = "<your api secret here>"
	userAgent = SDKName + "/" + SDKVersion
)

var ot = New(apiKey, apiSecret)

func TestNew(t *testing.T) {
	expect := &OpenTok{apiKey, apiSecret, defaultAPIHost, userAgent, http.DefaultClient}

	actual := New(apiKey, apiSecret)

	assert.Equal(t, expect, actual)
}

func TestOpenTok_SetAPIHost(t *testing.T) {
	newAPIHost := "http://example.com"

	err := ot.SetAPIHost(newAPIHost)

	assert.Nil(t, err)
	assert.Equal(t, newAPIHost, ot.apiHost)
}

func TestOpenTok_SetHttpClient(t *testing.T) {
	httpClient := &http.Client{
		Timeout: 120 * time.Second,
	}

	ot.SetHttpClient(httpClient)

	assert.Equal(t, httpClient, ot.httpClient)
}

func TestOpenTok_JwtToken(t *testing.T) {
	ot := New(apiKey, apiSecret)

	// Validate  project token
	tokenString, err := ot.jwtToken(projectToken)
	if err != nil {
		log.Fatal(err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		assert.IsType(t, &jwt.SigningMethodHMAC{}, token.Method)

		return []byte(apiSecret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	if assert.True(t, ok) {
		assert.Equal(t, "project", claims["ist"])
	}

	// Validate account token
	tokenString, err = ot.jwtToken(accountToken)
	if err != nil {
		log.Fatal(err)
	}

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		assert.IsType(t, &jwt.SigningMethodHMAC{}, token.Method)

		return []byte(apiSecret), nil
	})

	claims, ok = token.Claims.(jwt.MapClaims)

	if assert.True(t, ok) {
		assert.Equal(t, "account", claims["ist"])
	}
}
