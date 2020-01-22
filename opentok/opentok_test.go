package opentok

import (
	"log"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

const (
	apiKey    = "<your api key here>"
	apiSecret = "<your api secret here>"
)

var ot = New(apiKey, apiSecret)

func TestNew(t *testing.T) {
	expect := &OpenTok{apiKey, apiSecret, defaultAPIHost}

	actual := New(apiKey, apiSecret)

	assert.Equal(t, expect, actual)
}

func TestOpenTok_SetAPIHost(t *testing.T) {
	newAPIHost := "http://example.com"

	err := ot.SetAPIHost(newAPIHost)

	assert.Nil(t, err)
	assert.Equal(t, newAPIHost, ot.apiHost)
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
