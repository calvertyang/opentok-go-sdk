package opentok

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

const (
	apiKey    = "<your api key here>"
	apiSecret = "<your api secret here>"
)

var ot = New(apiKey, apiSecret)

func TestNew(t *testing.T) {
	expect := &OpenTok{apiKey, apiSecret, defaultAPIHost, false, http.DefaultClient}

	actual := New(apiKey, apiSecret)

	assert.Equal(t, expect, actual)
}

func TestOpenTok_SetAPIHost(t *testing.T) {
	newAPIHost := "http://example.com"

	err := ot.SetAPIHost(newAPIHost)

	assert.Nil(t, err)
	assert.Equal(t, newAPIHost, ot.apiHost)
}

func TestOpenTok_Debug(t *testing.T) {
	ot.Debug()
	assert.Equal(t, true, ot.debug)
	ot.debug = false
}

func TestOpenTok_SetHTTPClient(t *testing.T) {
	httpClient := &http.Client{
		Timeout: 120 * time.Second,
	}

	ot.SetHTTPClient(httpClient)

	assert.Equal(t, httpClient, ot.httpClient)
}

func TestOpenTok_GenAccountJWT(t *testing.T) {
	tokenString, err := ot.genProjectJWT()
	assert.Nil(t, err)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		assert.IsType(t, &jwt.SigningMethodHMAC{}, token.Method)

		return []byte(apiSecret), nil
	})
	assert.Nil(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)

	if assert.True(t, ok) {
		assert.Equal(t, string(projectToken), claims["ist"])
	}
}

func TestOpenTok_GenProjectJWT(t *testing.T) {
	tokenString, err := ot.genAccountJWT()
	assert.Nil(t, err)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		assert.IsType(t, &jwt.SigningMethodHMAC{}, token.Method)

		return []byte(apiSecret), nil
	})
	assert.Nil(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)

	if assert.True(t, ok) {
		assert.Equal(t, string(accountToken), claims["ist"])
	}
}

func TestOpenTok_SendRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	req, _ := http.NewRequest(http.MethodGet, ot.apiHost, nil)
	res, err := ot.sendRequest(context.Background(), req)

	assert.Nil(t, err)
	assert.IsType(t, &http.Response{}, res)
}

func TestOpenTok_ParseErrorResponse(t *testing.T) {
	sampleRes := http.Response{
		StatusCode: 400,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
			{
				"code": 10101,
				"message": "Invalid session id format",
				"description": "Invalid session id format"
			}
		`)),
	}

	expect := &ResponseError{
		StatusCode: 400,
		Message:    "Invalid session id format",
	}

	actual := parseErrorResponse(&sampleRes)

	assert.Equal(t, expect.Error(), actual.Error())
}
