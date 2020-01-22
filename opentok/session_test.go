package opentok

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOpenTok_CreateSession(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			[{
				"properties": null,
				"session_id": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"project_id": "40000001",
				"partner_id": "40000001",
				"create_dt": "Wed Jan 01 00:00:00 PST 2020",
				"session_status": null,
				"status_invalid": null,
				"media_server_hostname": null,
				"messaging_server_url": null,
				"messaging_url": null,
				"symphony_address": null,
				"ice_server": null,
				"session_segment_id": "29063f17-8d5f-4db1-9d96-93cd08c06619",
				"ice_servers": null,
				"ice_credential_expiration": 86100
			}]
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Session{
		SessionID:      "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		ProjectID:      "40000001",
		CreateDt:       "Wed Jan 01 00:00:00 PST 2020",
		MediaServerURL: "",
		OpenTok:        ot,
	}

	actual, err := ot.CreateSession(&SessionOptions{
		ArchiveMode: AutoArchived,
		MediaMode:   Routed,
	})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_GenerateToken(t *testing.T) {
	_, err := ot.GenerateToken("1_MX48eW91ciBhcGkga2V5IGhlcmU-fn4xNTc3ODY1NjAwMDAwfng3YjQ4TVFnRGYrWVFGcVBRaDh2VmZPS34", &TokenOptions{
		Role: Publisher,
	})

	assert.Nil(t, err)
}

func TestSession_GenerateToken(t *testing.T) {
	session := &Session{
		SessionID: "1_MX48eW91ciBhcGkga2V5IGhlcmU-fn4xNTc3ODY1NjAwMDAwfng3YjQ4TVFnRGYrWVFGcVBRaDh2VmZPS34",
		ProjectID: apiKey,
		OpenTok:   ot,
	}
	_, err := session.GenerateToken(&TokenOptions{
		Role: Publisher,
	})

	assert.Nil(t, err)
}

func TestDecodeSessionID(t *testing.T) {
	expect := &SessionIDInfo{
		APIKey:     "40000001",
		Location:   "",
		CreateTime: time.Unix(1577865600, 0),
	}

	actual, err := decodeSessionID("1_MX40MDAwMDAwMX5-MTU3Nzg2NTYwMDAwMH54N2I0OE1RZ0RmK1lRRnFQUWg4dlZmT0t-QX4")

	assert.Nil(t, err)

	if assert.NotEmpty(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestEncodeToken(t *testing.T) {
	tokenData := map[string]string{
		"session_id":                "1_MX40MDAwMDAwMX5-MTU3Nzg2NTYwMDAwMH54N2I0OE1RZ0RmK1lRRnFQUWg4dlZmT0t-QX4",
		"create_time":               "1577865600",
		"expire_time":               "1577865600",
		"nonce":                     "0.49893371771268225",
		"role":                      "publisher",
		"connection_data":           "foo=bar",
		"initial_layout_class_list": "",
	}

	expect := "T1==cGFydG5lcl9pZD08eW91ciBhcGkga2V5IGhlcmU+JnNpZz0yYjQyMzlkNjU4YTVkYmE0NGRhMGMyMmUzOTA2MWM5ZWI1ODQ1MTE1OmNvbm5lY3Rpb25fZGF0YT1mb28lM0RiYXImY3JlYXRlX3RpbWU9MTU3Nzg2NTYwMCZleHBpcmVfdGltZT0xNTc3ODY1NjAwJmluaXRpYWxfbGF5b3V0X2NsYXNzX2xpc3Q9Jm5vbmNlPTAuNDk4OTMzNzE3NzEyNjgyMjUmcm9sZT1wdWJsaXNoZXImc2Vzc2lvbl9pZD0xX01YNDBNREF3TURBd01YNS1NVFUzTnpnMk5UWXdNREF3TUg1NE4ySTBPRTFSWjBSbUsxbFJSbkZRVVdnNGRsWm1UMHQtUVg0"

	actual, err := encodeToken(tokenData, ot)

	assert.Nil(t, err)

	if assert.NotEmpty(t, actual) {
		assert.Equal(t, expect, actual)
	}
}
