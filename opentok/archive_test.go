package opentok

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenTok_StartArchive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
				"status": "started",
				"name": "example",
				"reason": "",
				"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"projectId": 40000001,
				"createdAt": 1579163008000,
				"size": 0,
				"duration": 0,
				"outputMode": "composed",
				"hasAudio": true,
				"hasVideo": true,
				"resolution": "1280x720",
				"url":null
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Archive{
		CreatedAt:  1579163008000,
		Duration:   0,
		HasAudio:   true,
		HasVideo:   true,
		ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
		Name:       "example",
		OutputMode: "composed",
		ProjectID:  40000001,
		Reason:     "",
		Resolution: "1280x720",
		SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		Size:       0,
		Status:     "started",
		URL:        nil,
		OpenTok:    ot,
	}

	actual, err := ot.StartArchive("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", &ArchiveOptions{
		Name: "example",
		Layout: &Layout{
			Type: PIP,
		},
		Resolution: HD,
	})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_StopArchive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
				"status": "stopped",
				"name": "example",
				"reason": "user initiated",
				"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"projectId": 40000001,
				"createdAt": 1579163008000,
				"size": 0,
				"duration": 0,
				"outputMode": "composed",
				"hasAudio": true,
				"hasVideo": true,
				"resolution": "1280x720",
				"url":null
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Archive{
		CreatedAt:  1579163008000,
		Duration:   0,
		HasAudio:   true,
		HasVideo:   true,
		ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
		Name:       "example",
		OutputMode: "composed",
		ProjectID:  40000001,
		Reason:     "user initiated",
		Resolution: "1280x720",
		SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		Size:       0,
		Status:     "stopped",
		URL:        nil,
		OpenTok:    ot,
	}

	actual, err := ot.StopArchive("c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea")

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_ListArchives(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"count": 1,
				"items": [{
					"id": "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
					"status": "uploaded",
					"name": "example",
					"reason": "user initiated",
					"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
					"projectId": 40000001,
					"createdAt": 1579163008000,
					"size": 355816,
					"duration": 34,
					"outputMode": "composed",
					"hasAudio": true,
					"hasVideo": true,
					"resolution": "1280x720",
					"url":null
				}]
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &ArchiveList{
		Count: 1,
		Items: []*Archive{
			{
				CreatedAt:  1579163008000,
				Duration:   34,
				HasAudio:   true,
				HasVideo:   true,
				ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
				Name:       "example",
				OutputMode: "composed",
				ProjectID:  40000001,
				Reason:     "user initiated",
				Resolution: "1280x720",
				SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				Size:       355816,
				Status:     "uploaded",
				URL:        nil,
				OpenTok:    ot,
			},
		},
	}

	actual, err := ot.ListArchives(&ArchiveListOptions{})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_GetArchive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
				"status": "uploaded",
				"name": "example",
				"reason": "user initiated",
				"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"projectId": 40000001,
				"createdAt": 1579163008000,
				"size": 355816,
				"duration": 34,
				"outputMode": "composed",
				"hasAudio": true,
				"hasVideo": true,
				"resolution": "1280x720",
				"url":null
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Archive{
		CreatedAt:  1579163008000,
		Duration:   34,
		HasAudio:   true,
		HasVideo:   true,
		ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
		Name:       "example",
		OutputMode: "composed",
		ProjectID:  40000001,
		Reason:     "user initiated",
		Resolution: "1280x720",
		SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		Size:       355816,
		Status:     "uploaded",
		URL:        nil,
		OpenTok:    ot,
	}

	actual, err := ot.GetArchive("c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea")

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_DeleteArchive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	err := ot.DeleteArchive("c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea")

	assert.Nil(t, err)
}

func TestOpenTok_SetArchiveStorage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"type": "s3",
				"config": {
					"accessKey": "myUsername",
					"bucket": "example",
					"secretKey": "myPassword"
				},
				"fallback": "none"
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &StorageOptions{
		Type: S3,
		Config: map[string]interface{}{
			"accessKey": "myUsername",
			"bucket":    "example",
			"secretKey": "myPassword",
		},
		Fallback: "none",
	}

	actual, err := ot.SetArchiveStorage(&StorageOptions{
		Type: S3,
		Config: &AmazonS3Config{
			AccessKey: "myUsername",
			SecretKey: "myPassword",
			Bucket:    "example",
		},
		Fallback: "none",
	})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_DeleteArchiveStorage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	err := ot.DeleteArchiveStorage()

	assert.Nil(t, err)
}

func TestOpenTok_SetArchiveLayout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
				"status": "started",
				"name": "example",
				"reason": "",
				"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"projectId": 40000001,
				"createdAt": 1579163008000,
				"size": 0,
				"duration": 0,
				"outputMode": "composed",
				"hasAudio": true,
				"hasVideo": true,
				"resolution": "1280x720",
				"url":null
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Archive{
		CreatedAt:  1579163008000,
		Duration:   0,
		HasAudio:   true,
		HasVideo:   true,
		ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
		Name:       "example",
		OutputMode: "composed",
		ProjectID:  40000001,
		Reason:     "",
		Resolution: "1280x720",
		SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		Size:       0,
		Status:     "started",
		URL:        nil,
		OpenTok:    ot,
	}

	actual, err := ot.SetArchiveLayout("c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea", &Layout{
		Type: BestFit,
	})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestArchive_Stop(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
				"status": "stopped",
				"name": "example",
				"reason": "user initiated",
				"sessionId": "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
				"projectId": 40000001,
				"createdAt": 1579163008000,
				"size": 0,
				"duration": 0,
				"outputMode": "composed",
				"hasAudio": true,
				"hasVideo": true,
				"resolution": "1280x720",
				"url":null
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Archive{
		CreatedAt:  1579163008000,
		Duration:   0,
		HasAudio:   true,
		HasVideo:   true,
		ID:         "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
		Name:       "example",
		OutputMode: "composed",
		ProjectID:  40000001,
		Reason:     "user initiated",
		Resolution: "1280x720",
		SessionID:  "1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4",
		Size:       0,
		Status:     "stopped",
		URL:        nil,
		OpenTok:    ot,
	}

	archive := &Archive{
		ID:      "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
		OpenTok: ot,
	}
	actual, err := archive.Stop()

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestArchive_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	archive := &Archive{
		ID:      "c9c87fbb-f91b-49bf-a6b4-fd0dbe16caea",
		OpenTok: ot,
	}
	err := archive.Delete()

	assert.Nil(t, err)
}
