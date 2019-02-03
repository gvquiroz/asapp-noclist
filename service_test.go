package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthWithValidResponse(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == `/auth` {
			w.Header().Set("Badsec-Authentication-Token", "token")
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	c := &BADSECClient{
		Client:         &http.Client{},
		BASDECEndpoint: server.URL,
	}

	token, err := c.getAuthToken()

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, "token", token)
}

func TestAuthWithInvalidServiceUnavaialable(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == `/auth` {
			w.Header().Set("Badsec-Authentication-Token", "token")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	defer server.Close()

	c := &BADSECClient{
		Client:         &http.Client{},
		BASDECEndpoint: server.URL,
	}

	_, err := c.getAuthToken()

	if err == nil {
		assert.Fail(t, "getAuthToken must get an error if service is Unavailable")
	}

	assert.Contains(t, err.Error(), "Invalid response from API")

}
func TestAuthWithoutAuthHeader(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == `/auth` {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	c := &BADSECClient{
		Client:         &http.Client{},
		BASDECEndpoint: server.URL,
	}

	_, err := c.getAuthToken()

	if err == nil {
		assert.Fail(t, "getAuthToken must get an error if Auth Header is not present")
	}

	assert.Contains(t, err.Error(), "No BADSEC auth in response")

}
func TestAuthWithInvalidHeaderValue(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == `/auth` {
			w.Header().Set("Badsec-Authentication-Token", "")
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	c := &BADSECClient{
		Client:         &http.Client{},
		BASDECEndpoint: server.URL,
	}

	_, err := c.getAuthToken()

	if err == nil {
		assert.Fail(t, "getAuthToken must get an error if Auth Header is invalid")
	}

	assert.Contains(t, err.Error(), "No BADSEC auth in response")

}

func TestGenerateChecksumUsingURLAndToken(t *testing.T) {
	result := generateChecksum("/user", "aAmazingToken")

	// Calculated using https://www.xorbin.com/tools/sha256-hash-calculator
	manuallyCalculatedResult := "6222ee011cfc9ee987cdaa10a7ecf61559442d5670d67c2db01a9d82b406cb1e"

	assert.Equal(t, manuallyCalculatedResult, result)

}

func TestNewBADSECClient(t *testing.T) {

	tockenMock := "JustAnotherMockToken"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == `/auth` {
			w.Header().Set("Badsec-Authentication-Token", tockenMock)
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := NewService(server.URL)

	userChecksumResult := generateChecksum("/users", tockenMock)

	assert.Equal(t, client.BASDECEndpoint, server.URL)
	assert.Equal(t, client.UsersChecksum, userChecksumResult)

}
