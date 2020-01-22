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
