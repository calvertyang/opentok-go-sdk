package opentok

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenTok_SendSessionSignal(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	err := ot.SendSessionSignal("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", &SignalData{
		Type: "foo",
		Data: "bar",
	})

	assert.Nil(t, err)
}

func TestOpenTok_SendConnectionSignal(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	err := ot.SendConnectionSignal("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", "18145975-97c8-4802-8975-fc8408d56d5c", &SignalData{
		Type: "foo",
		Data: "bar",
	})

	assert.Nil(t, err)
}

func TestSession_SendSignal(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	session := &Session{
		SessionID: "1_MX48eW91ciBhcGkga2V5IGhlcmU-fn4xNTc3ODY1NjAwMDAwfng3YjQ4TVFnRGYrWVFGcVBRaDh2VmZPS34",
		ProjectID: apiKey,
		OpenTok:   ot,
	}

	err := session.SendSignal(&SignalData{
		Type: "foo",
		Data: "bar",
	})

	assert.Nil(t, err)
}

func TestSession_SendConnectionSignal(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	session := &Session{
		SessionID: "1_MX48eW91ciBhcGkga2V5IGhlcmU-fn4xNTc3ODY1NjAwMDAwfng3YjQ4TVFnRGYrWVFGcVBRaDh2VmZPS34",
		ProjectID: apiKey,
		OpenTok:   ot,
	}

	err := session.SendConnectionSignal("18145975-97c8-4802-8975-fc8408d56d5c", &SignalData{
		Type: "foo",
		Data: "bar",
	})

	assert.Nil(t, err)
}
