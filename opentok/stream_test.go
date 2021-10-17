package opentok

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenTok_ListStreams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"count": 1,
				"items": [{
					"id": "d962b966-964d-4f18-be3f-e0b181a43b0e",
					"videoType": "camera",
					"name": "",
					"layoutClassList": []
				}]
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &StreamList{
		Count: 1,
		Items: []*Stream{
			{
				ID:              "d962b966-964d-4f18-be3f-e0b181a43b0e",
				VideoType:       "camera",
				Name:            "",
				LayoutClassList: []string{},
			},
		},
	}

	actual, err := ot.ListStreams("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4")

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_GetStream(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "d962b966-964d-4f18-be3f-e0b181a43b0e",
				"videoType": "camera",
				"name": "",
				"layoutClassList": []
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Stream{
		ID:              "d962b966-964d-4f18-be3f-e0b181a43b0e",
		VideoType:       "camera",
		Name:            "",
		LayoutClassList: []string{},
	}

	actual, err := ot.GetStream("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", "d962b966-964d-4f18-be3f-e0b181a43b0e")

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_SetStreamClassLists(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"count": 1,
				"items": [{
					"id": "d962b966-964d-4f18-be3f-e0b181a43b0e",
					"videoType": "camera",
					"name": "",
					"layoutClassList": ["full"]
				}]
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &StreamList{
		Count: 1,
		Items: []*Stream{
			{
				ID:              "d962b966-964d-4f18-be3f-e0b181a43b0e",
				VideoType:       "camera",
				Name:            "",
				LayoutClassList: []string{"full"},
			},
		},
	}

	actual, err := ot.SetStreamClassLists("1_QX90NjQ2MCY0Nm6-MTU4QTO4NzE5NTkyOX4yUy2OZndKQExJR0NyalcvNktmTzBpSnp-QX4", &StreamClassOptions{
		Items: []*StreamClass{
			{
				ID:              "d962b966-964d-4f18-be3f-e0b181a43b0e",
				LayoutClassList: []string{"full"},
			},
		},
	})

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}
