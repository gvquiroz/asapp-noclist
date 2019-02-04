package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// BADSECClient Represents a interface to communicate with BADSEC API
type BADSECClient struct {
	UsersChecksum  string
	Client         *http.Client
	BASDECEndpoint string
}

// getUsers receive the number of times that will try to attemp the request if the firstone fails. Number of request will be 1 + retryAttemps
func (b BADSECClient) getUsers(retryAttemps int) (string, error) {

	req, err := http.NewRequest("GET", b.BASDECEndpoint+"/users", nil)
	req.Header.Add("X-Request-Checksum", b.UsersChecksum)
	response, err := b.Client.Do(req)

	var usersJSON string

	if err != nil {
		return usersJSON, err
	}

	// Any response code other than 200 is marked as a failure
	if response.StatusCode != http.StatusOK {
		if retryAttemps > 0 {
			log.Println("[ Retrying ] Invalid response from API - Attemp number", retryAttemps)
			return b.getUsers(retryAttemps - 1)
		}

		return usersJSON, errors.New("Invalid response from API ")
	}

	defer response.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(response.Body)

	s := strings.Split(string(bodyBytes), "\n")
	u, _ := json.Marshal(s)

	usersJSON = string(u)

	return usersJSON, nil
}

// NewService returns a BADSECClient to interact with BADSEC API
func NewService(APIEndpoint string) *BADSECClient {

	retryAttepms := 3

	c := &BADSECClient{
		Client:         &http.Client{},
		BASDECEndpoint: APIEndpoint,
	}

	token, err := c.getAuthToken(retryAttepms)

	if err != nil {
		log.Fatal(err)
	}

	c.UsersChecksum = generateChecksum("/users", token)

	return c

}

// getAuthToken receive the number of times that will try to attemp the request if the firstone fails. Number of request will be 1 + retryAttemps
func (b BADSECClient) getAuthToken(retryAttemps int) (string, error) {
	response, err := b.Client.Get(b.BASDECEndpoint + "/auth")
	var token string

	if err != nil {
		return token, err
	}

	// Any response code other than 200 is marked as a failure
	if response.StatusCode != http.StatusOK {

		if retryAttemps > 0 {
			log.Println("[ Retrying ] Invalid response from API - Attemp number", retryAttemps)
			return b.getAuthToken(retryAttemps - 1)
		}

		return token, errors.New("Invalid response from API ")

	}

	// We check if the header is present
	if tokenValue, ok := response.Header["Badsec-Authentication-Token"]; ok {
		token = tokenValue[0]
	}

	// We check if the value of the header is not empty
	if token == "" {
		return token, errors.New("No BADSEC auth in response")
	}

	return token, nil
}

func generateChecksum(path, authToken string) string {
	sha256 := sha256.Sum256([]byte(authToken + path))
	return fmt.Sprintf("%x", sha256)
}
