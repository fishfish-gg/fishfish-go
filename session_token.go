package fishfish

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

type SessionToken struct {
	token string
	mx    sync.Mutex
}

type createTokenRequest struct {
	Permissions []string `json:"Permissions"`
}

type TokenResponse struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}

func (client *Client) fetchSessionToken() (string, error) {
	var err error
	var token string

	requestBody := createTokenRequest{
		Permissions: client.config.Permissions,
	}

	requestBodyJson, err := json.Marshal(requestBody)
	apiURL := client.getAPIUrl("users/@me/tokens")
	reqBody := bytes.NewBuffer(requestBodyJson)

	req, _ := http.NewRequest("POST", apiURL, reqBody)
	req.Header.Set("Authorization", client.config.Auth)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return token, err
	}

	var tokenResponse TokenResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &tokenResponse)

	token = tokenResponse.Token

	return token, err
}

func (client *Client) updateSessionToken() string {
	sessionToken, err := client.fetchSessionToken()

	if err != nil {
		log.Fatal(err)
	}

	client.sessionToken.mx.Lock()
	client.sessionToken.token = sessionToken

	client.sessionToken.mx.Unlock()

	return sessionToken
}

func (client *Client) getSessionToken() string {
	client.sessionToken.mx.Lock()

	token := client.sessionToken.token

	defer client.sessionToken.mx.Unlock()

	return token
}
