package opentok

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenTok_StartBroadcast(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "ce872e0d-4997-440a-a0a5-10ce715b54cf",
				"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"projectId": 40000001,
				"createdAt": 1579163008000,
				"updatedAt": 1579163008000,
				"resolution": "1280x720",
				"status": "started",
				"broadcastUrls": {
					"hls": "",
					"rtmp": [{
						"id": "foo",
						"status": "connecting",
						"serverUrl": "rtmps://myfooserver/myfooapp",
						"streamName": "myfoostream"
					}, {
						"id": "bar",
						"status": "connecting",
						"serverUrl": "rtmp://mybarserver/mybarapp",
						"streamName": "mybarstream"
					}]
				}
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Broadcast{
		ID:         "ce872e0d-4997-440a-a0a5-10ce715b54cf",
		SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		ProjectID:  40000001,
		CreatedAt:  1579163008000,
		UpdatedAt:  1579163008000,
		Resolution: "1280x720",
		Status:     "started",
		BroadcastURLs: &BroadcastURLs{
			HLS: "",
			RTMP: []*RTMPConfig{
				{
					ID:         "foo",
					Status:     "connecting",
					ServerURL:  "rtmps://myfooserver/myfooapp",
					StreamName: "myfoostream",
				},
				{
					ID:         "bar",
					Status:     "connecting",
					ServerURL:  "rtmp://mybarserver/mybarapp",
					StreamName: "mybarstream",
				},
			},
		},
		OpenTok: ot,
	}

	actual, err := ot.StartBroadcast("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", &BroadcastOptions{
		Layout: &Layout{
			Type: PIP,
		},
		MaxDuration: 60,
		Outputs: &BroadcastOutputOptions{
			RTMP: []*RTMPConfig{
				{
					ID:         "foo",
					ServerURL:  "rtmps://myfooserver/myfooapp",
					StreamName: "myfoostream",
				},
				{
					ID:         "bar",
					ServerURL:  "rtmp://mybarserver/mybarapp",
					StreamName: "mybarstream",
				},
			},
		},
		Resolution: HD,
	})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_StopBroadcast(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "ce872e0d-4997-440a-a0a5-10ce715b54cf",
				"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"projectId": 40000001,
				"createdAt": 1579163008000,
				"updatedAt": 1579163009000,
				"resolution": "1280x720",
				"status": "stopped",
				"broadcastUrls": null
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Broadcast{
		ID:            "ce872e0d-4997-440a-a0a5-10ce715b54cf",
		SessionID:     "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		ProjectID:     40000001,
		CreatedAt:     1579163008000,
		UpdatedAt:     1579163009000,
		Resolution:    "1280x720",
		Status:        "stopped",
		BroadcastURLs: nil,
		OpenTok:       ot,
	}

	actual, err := ot.StopBroadcast("ce872e0d-4997-440a-a0a5-10ce715b54cf")

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_ListBroadcasts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"count": 1,
				"items": [{
					"id": "ce872e0d-4997-440a-a0a5-10ce715b54cf",
					"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
					"projectId": 40000001,
					"createdAt": 1579163008000,
					"updatedAt": 1579163008000,
					"resolution": "1280x720",
					"status": "started",
					"broadcastUrls": {
						"hls": "",
						"rtmp": [{
							"id": "foo",
							"status": "connecting",
							"serverUrl": "rtmps://myfooserver/myfooapp",
							"streamName": "myfoostream"
						}, {
							"id": "bar",
							"status": "connecting",
							"serverUrl": "rtmp://mybarserver/mybarapp",
							"streamName": "mybarstream"
						}]
					}
				}]
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &BroadcastList{
		Count: 1,
		Items: []*Broadcast{
			{
				ID:         "ce872e0d-4997-440a-a0a5-10ce715b54cf",
				SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				ProjectID:  40000001,
				CreatedAt:  1579163008000,
				UpdatedAt:  1579163008000,
				Resolution: "1280x720",
				Status:     "started",
				BroadcastURLs: &BroadcastURLs{
					HLS: "",
					RTMP: []*RTMPConfig{
						{
							ID:         "foo",
							Status:     "connecting",
							ServerURL:  "rtmps://myfooserver/myfooapp",
							StreamName: "myfoostream",
						},
						{
							ID:         "bar",
							Status:     "connecting",
							ServerURL:  "rtmp://mybarserver/mybarapp",
							StreamName: "mybarstream",
						},
					},
				},
				OpenTok: ot,
			},
		},
	}

	actual, err := ot.ListBroadcasts(&BroadcastListOptions{})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_GetBroadcast(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "ce872e0d-4997-440a-a0a5-10ce715b54cf",
				"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"projectId": 40000001,
				"createdAt": 1579163008000,
				"updatedAt": 1579163009000,
				"resolution": "1280x720",
				"status": "stopped",
				"broadcastUrls": null
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Broadcast{
		ID:            "ce872e0d-4997-440a-a0a5-10ce715b54cf",
		SessionID:     "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		ProjectID:     40000001,
		CreatedAt:     1579163008000,
		UpdatedAt:     1579163009000,
		Resolution:    "1280x720",
		Status:        "stopped",
		BroadcastURLs: nil,
		OpenTok:       ot,
	}

	actual, err := ot.GetBroadcast("ce872e0d-4997-440a-a0a5-10ce715b54cf")

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_SetBroadcastLayout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "ce872e0d-4997-440a-a0a5-10ce715b54cf",
				"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"projectId": 40000001,
				"createdAt": 1579163008000,
				"updatedAt": 1579163008000,
				"resolution": "1280x720",
				"status": "started",
				"broadcastUrls": {
					"hls": "",
					"rtmp": [{
						"id": "foo",
						"status": "connecting",
						"serverUrl": "rtmps://myfooserver/myfooapp",
						"streamName": "myfoostream"
					}, {
						"id": "bar",
						"status": "connecting",
						"serverUrl": "rtmp://mybarserver/mybarapp",
						"streamName": "mybarstream"
					}]
				}
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Broadcast{
		ID:         "ce872e0d-4997-440a-a0a5-10ce715b54cf",
		SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		ProjectID:  40000001,
		CreatedAt:  1579163008000,
		UpdatedAt:  1579163008000,
		Resolution: "1280x720",
		Status:     "started",
		BroadcastURLs: &BroadcastURLs{
			HLS: "",
			RTMP: []*RTMPConfig{
				{
					ID:         "foo",
					Status:     "connecting",
					ServerURL:  "rtmps://myfooserver/myfooapp",
					StreamName: "myfoostream",
				},
				{
					ID:         "bar",
					Status:     "connecting",
					ServerURL:  "rtmp://mybarserver/mybarapp",
					StreamName: "mybarstream",
				},
			},
		},
		OpenTok: ot,
	}

	actual, err := ot.SetBroadcastLayout("ce872e0d-4997-440a-a0a5-10ce715b54cf", &Layout{
		Type: BestFit,
	})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestBroadcast_Stop(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "ce872e0d-4997-440a-a0a5-10ce715b54cf",
				"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"projectId": 40000001,
				"createdAt": 1579163008000,
				"updatedAt": 1579163009000,
				"resolution": "1280x720",
				"status": "stopped",
				"broadcastUrls": null
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Broadcast{
		ID:            "ce872e0d-4997-440a-a0a5-10ce715b54cf",
		SessionID:     "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		ProjectID:     40000001,
		CreatedAt:     1579163008000,
		UpdatedAt:     1579163009000,
		Resolution:    "1280x720",
		Status:        "stopped",
		BroadcastURLs: nil,
		OpenTok:       ot,
	}

	broadcast := &Broadcast{
		ID:      "ce872e0d-4997-440a-a0a5-10ce715b54cf",
		OpenTok: ot,
	}
	actual, err := broadcast.Stop()

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}
