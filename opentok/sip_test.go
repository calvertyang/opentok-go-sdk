package opentok

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenTok_Dial(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "ab31819a-cd52-4da4-8b3e-fb9803337c17",
				"connectionId": "3a6aa409-bfc5-462c-a9c7-59b72aeebf69",
				"streamId": "f1a58131-7b2c-4fa8-b2a7-64fdc6b2c0f6"
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &SIPCall{
		ID:           "ab31819a-cd52-4da4-8b3e-fb9803337c17",
		ConnectionID: "3a6aa409-bfc5-462c-a9c7-59b72aeebf69",
		StreamID:     "f1a58131-7b2c-4fa8-b2a7-64fdc6b2c0f6",
	}

	actual, err := ot.Dial("1_MX48eW91ciBhcGkga2V5IGhlcmU-fn4xNTc3ODY1NjAwMDAwfng3YjQ4TVFnRGYrWVFGcVBRaDh2VmZPS34", &DialOptions{
		SIP: &SIP{
			URI: "sip:user@sip.example.comwhen;transport=tls",
			Headers: &SIPHeaders{
				"X-Foo": "bar",
			},
			Auth: &SIPAuth{
				Username: "username",
				Password: "password",
			},
			Secure: true,
		},
	})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_SendDTMF(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	err := ot.SendDTMF("1_MX48eW91ciBhcGkga2V5IGhlcmU-fn4xNTc3ODY1NjAwMDAwfng3YjQ4TVFnRGYrWVFGcVBRaDh2VmZPS34", "1713")

	assert.Nil(t, err)
}

func TestOpenTok_SendDTMFToClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	err := ot.SendDTMFToClient("1_MX48eW91ciBhcGkga2V5IGhlcmU-fn4xNTc3ODY1NjAwMDAwfng3YjQ4TVFnRGYrWVFGcVBRaDh2VmZPS34", "3a6aa409-bfc5-462c-a9c7-59b72aeebf69", "1713")

	assert.Nil(t, err)
}

func TestSession_Dial(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "ab31819a-cd52-4da4-8b3e-fb9803337c17",
				"connectionId": "3a6aa409-bfc5-462c-a9c7-59b72aeebf69",
				"streamId": "f1a58131-7b2c-4fa8-b2a7-64fdc6b2c0f6"
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &SIPCall{
		ID:           "ab31819a-cd52-4da4-8b3e-fb9803337c17",
		ConnectionID: "3a6aa409-bfc5-462c-a9c7-59b72aeebf69",
		StreamID:     "f1a58131-7b2c-4fa8-b2a7-64fdc6b2c0f6",
	}

	session := &Session{
		SessionID: "1_MX48eW91ciBhcGkga2V5IGhlcmU-fn4xNTc3ODY1NjAwMDAwfng3YjQ4TVFnRGYrWVFGcVBRaDh2VmZPS34",
		ProjectID: apiKey,
		OpenTok:   ot,
	}

	actual, err := session.Dial(&DialOptions{
		SIP: &SIP{
			URI: "sip:user@sip.example.comwhen;transport=tls",
			Headers: &SIPHeaders{
				"X-Foo": "bar",
			},
			Auth: &SIPAuth{
				Username: "username",
				Password: "password",
			},
			Secure: true,
		},
	})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestSession_SendDTMF(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	session := &Session{
		SessionID: "1_MX48eW91ciBhcGkga2V5IGhlcmU-fn4xNTc3ODY1NjAwMDAwfng3YjQ4TVFnRGYrWVFGcVBRaDh2VmZPS34",
		ProjectID: apiKey,
		OpenTok:   ot,
	}

	err := session.SendDTMF("1713")

	assert.Nil(t, err)
}

func TestSession_SendDTMFToClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	session := &Session{
		SessionID: "1_MX48eW91ciBhcGkga2V5IGhlcmU-fn4xNTc3ODY1NjAwMDAwfng3YjQ4TVFnRGYrWVFGcVBRaDh2VmZPS34",
		ProjectID: apiKey,
		OpenTok:   ot,
	}

	err := session.SendDTMFToClient("3a6aa409-bfc5-462c-a9c7-59b72aeebf69", "1713")

	assert.Nil(t, err)
}
