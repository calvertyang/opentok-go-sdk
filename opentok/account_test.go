package opentok

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenTok_CreateProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "40000001",
				"secret": "ba7816bf8f01cfea414140de5dae2223b00361a3",
				"status": "VALID",
				"name": "example",
				"createdAt": 1579163008000,
				"environmentName": "default",
				"environmentDescription": "Standard Environment"
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Project{
		ID:                     "40000001",
		Secret:                 "ba7816bf8f01cfea414140de5dae2223b00361a3",
		Status:                 "VALID",
		Name:                   "example",
		CreatedAt:              1579163008000,
		EnvironmentName:        "default",
		EnvironmentDescription: "Standard Environment",
	}

	actual, err := ot.CreateProject("example")

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_ListProjects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			[{
				"id": "40000001",
				"secret": "ba7816bf8f01cfea414140de5dae2223b00361a3",
				"status": "VALID",
				"name": "example",
				"createdAt": 1579163008000,
				"environmentName": "default",
				"environmentDescription": "Standard Environment"
			}, {
				"id": "40000002",
				"secret": "fcde2b2edba56bf408601fb721fe9b5c338d10ee",
				"status": "VALID",
				"name": "foo",
				"createdAt": 1579163506000,
				"environmentName": "default",
				"environmentDescription": "Standard Environment"
			}]
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := []*Project{
		&Project{
			ID:                     "40000001",
			Secret:                 "ba7816bf8f01cfea414140de5dae2223b00361a3",
			Status:                 "VALID",
			Name:                   "example",
			CreatedAt:              1579163008000,
			EnvironmentName:        "default",
			EnvironmentDescription: "Standard Environment",
		},
		&Project{
			ID:                     "40000002",
			Secret:                 "fcde2b2edba56bf408601fb721fe9b5c338d10ee",
			Status:                 "VALID",
			Name:                   "foo",
			CreatedAt:              1579163506000,
			EnvironmentName:        "default",
			EnvironmentDescription: "Standard Environment",
		},
	}

	actual, err := ot.ListProjects()

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_GetProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "40000001",
				"secret": "ba7816bf8f01cfea414140de5dae2223b00361a3",
				"status": "VALID",
				"name": "example",
				"createdAt": 1579163008000,
				"environmentName": "default",
				"environmentDescription": "Standard Environment"
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Project{
		ID:                     "40000001",
		Secret:                 "ba7816bf8f01cfea414140de5dae2223b00361a3",
		Status:                 "VALID",
		Name:                   "example",
		CreatedAt:              1579163008000,
		EnvironmentName:        "default",
		EnvironmentDescription: "Standard Environment",
	}

	actual, err := ot.GetProject("40000001")

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_ChangeProjectStatus_Active(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "40000001",
				"secret": "ba7816bf8f01cfea414140de5dae2223b00361a3",
				"status": "ACTIVE",
				"name": "example",
				"createdAt": 1579163008000,
				"environmentName": "default",
				"environmentDescription": "Standard Environment"
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Project{
		ID:                     "40000001",
		Secret:                 "ba7816bf8f01cfea414140de5dae2223b00361a3",
		Status:                 "ACTIVE",
		Name:                   "example",
		CreatedAt:              1579163008000,
		EnvironmentName:        "default",
		EnvironmentDescription: "Standard Environment",
	}

	actual, err := ot.ChangeProjectStatus("40000001", ProjectActive)

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_ChangeProjectStatus_Suspend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "40000001",
				"secret": "ba7816bf8f01cfea414140de5dae2223b00361a3",
				"status": "SUSPENDED",
				"name": "example",
				"createdAt": 1579163008000,
				"environmentName": "default",
				"environmentDescription": "Standard Environment"
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Project{
		ID:                     "40000001",
		Secret:                 "ba7816bf8f01cfea414140de5dae2223b00361a3",
		Status:                 "SUSPENDED",
		Name:                   "example",
		CreatedAt:              1579163008000,
		EnvironmentName:        "default",
		EnvironmentDescription: "Standard Environment",
	}

	actual, err := ot.ChangeProjectStatus("40000001", ProjectActive)

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_RefreshProjectSecret(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"id": "40000001",
				"secret": "6a1e361fc002c0b2b51a821d7ce11f96e1887aea",
				"status": "ACTIVE",
				"name": "example",
				"createdAt": 1579163008000,
				"environmentName": "default",
				"environmentDescription": "Standard Environment"
			}
		`))
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	expect := &Project{
		ID:                     "40000001",
		Secret:                 "6a1e361fc002c0b2b51a821d7ce11f96e1887aea",
		Status:                 "ACTIVE",
		Name:                   "example",
		CreatedAt:              1579163008000,
		EnvironmentName:        "default",
		EnvironmentDescription: "Standard Environment",
	}

	actual, err := ot.RefreshProjectSecret("40000001")

	assert.Nil(t, err)

	if assert.NotNil(t, actual) {
		assert.Equal(t, expect, actual)
	}
}

func TestOpenTok_DeleteProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	ot.SetAPIHost(ts.URL)

	err := ot.DeleteProject("40000001")

	assert.Nil(t, err)
}
