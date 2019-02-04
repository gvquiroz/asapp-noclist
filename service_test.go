package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersWithValidResponse(t *testing.T) {

	users := `9757263792576857988
	7789651288773276582
	1628388650278268240`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == `/users` {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, users)
		}
	}))
	defer server.Close()

	c := &BADSECClient{
		Client:         &http.Client{},
		BASDECEndpoint: server.URL,
		UsersChecksum:  "mockCheckSum",
	}

	numberOfRetryAttemps := 3

	users, err := c.getUsers(numberOfRetryAttemps)

	if err != nil {
		log.Fatal(err)
	}

	usersJSON := fmt.Sprintf(`["9757263792576857988","\t7789651288773276582","\t1628388650278268240"]`)

	assert.Equal(t, usersJSON, users)
}

func TestUsersWithInvalidServiceUnavaialable(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == `/users` {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	defer server.Close()

	c := &BADSECClient{
		Client:         &http.Client{},
		BASDECEndpoint: server.URL,
		UsersChecksum:  "mockCheckSum",
	}

	numberOfRetryAttemps := 2

	_, err := c.getUsers(numberOfRetryAttemps)

	if err == nil {
		assert.Fail(t, "getUsers must get an error if service is Unavailable")
	}

	assert.Contains(t, err.Error(), "Invalid response from API")

}
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

	numberOfRetryAttemps := 0
	token, _ := c.getAuthToken(numberOfRetryAttemps)

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

	_, err := c.getAuthToken(0)

	if err == nil {
		assert.Fail(t, "getAuthToken must get an error if service is Unavailable")
	}

	assert.Contains(t, err.Error(), "Invalid response from API")

}

func TestAuthWithConnectionRefusedError(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	// Simulates that the server is not running
	server.Close()

	c := &BADSECClient{
		Client:         &http.Client{},
		BASDECEndpoint: server.URL,
	}

	_, err := c.getAuthToken(0)

	if err == nil {
		assert.Fail(t, "getAuthToken must get an error if service is Unavailable")
	}

	assert.Contains(t, err.Error(), "connect: connection refused")

}

func TestUsersWithConnectionRefusedError(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	// Simulates that the server is not running
	server.Close()

	c := &BADSECClient{
		Client:         &http.Client{},
		BASDECEndpoint: server.URL,
		UsersChecksum:  "mockCheckSum",
	}

	numberOfRetryAttemps := 0
	_, err := c.getUsers(numberOfRetryAttemps)

	if err == nil {
		assert.Fail(t, "getUsers must get an error if service is Unavailable")
	}

	assert.Contains(t, err.Error(), "connection refused")

}
func TestAuthWithInvalidServiceUnavaialableAndRetryAttemps(t *testing.T) {

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

	numberOfRetryAttemps := 3
	_, err := c.getAuthToken(numberOfRetryAttemps)

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

	numberOfRetryAttemps := 0
	_, err := c.getAuthToken(numberOfRetryAttemps)

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

	numberOfRetryAttemps := 0
	_, err := c.getAuthToken(numberOfRetryAttemps)

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
