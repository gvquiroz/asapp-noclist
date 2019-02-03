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

type BADSECClient struct {
	UsersChecksum  string
	Client         *http.Client
	BASDECEndpoint string
}

func (b BADSECClient) getUsers() (string, error) {

	req, err := http.NewRequest("GET", b.BASDECEndpoint+"/users", nil)
	req.Header.Add("X-Request-Checksum", b.UsersChecksum)
	response, err := b.Client.Do(req)

	var usersJSON string

	if err != nil {
		return usersJSON, err
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return usersJSON, err
	}

	s := strings.Split(string(bodyBytes), "\n")
	u, err := json.Marshal(s)

	if err != nil {
		return usersJSON, err
	}

	usersJSON = string(u)

	return usersJSON, nil
}

func NewService(APIEndpoint string) *BADSECClient {

	c := &BADSECClient{
		Client:         &http.Client{},
		BASDECEndpoint: APIEndpoint,
	}

	token, err := c.getAuthToken()

	if err != nil {
		log.Fatal(err)
	}

	c.UsersChecksum = generateChecksum("/users", token)

	return c

}

func (b BADSECClient) getAuthToken() (string, error) {
	response, err := b.Client.Get(b.BASDECEndpoint + "/auth")
	var token string

	if err != nil {
		return token, err
	}

	// Any response code other than 200 is marked as a failure
	if response.StatusCode != http.StatusOK {
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
