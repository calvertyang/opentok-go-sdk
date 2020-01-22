package opentok

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenTok_ForceDisconnect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	err := ot.ForceDisconnect("40000001", "efdf2fc7-bd6e-4871-9c1d-531f7f6a9486")

	assert.Nil(t, err)
}
