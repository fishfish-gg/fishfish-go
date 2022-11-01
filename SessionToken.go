package fishfish

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (client *FishFishClient) fetchSessionToken() (string, error) {
	var err error
	var token string

	requestBody := createTokenRequest{
		Permissions: client.config.Permissions,
	}

	requestBodyJson, err := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", client.getAPIUrl("users/@me/tokens"), bytes.NewBuffer(requestBodyJson))
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

func (client *FishFishClient) updateSessionToken() string {
	sessionToken, err := client.fetchSessionToken()

	if err != nil {
		log.Fatal(err)
	}

	client.sessionToken.mx.Lock()
	client.sessionToken.token = sessionToken

	client.sessionToken.mx.Unlock()

	return sessionToken
}

func (client *FishFishClient) getSessionToken() string {
	client.sessionToken.mx.Lock()

	token := client.sessionToken.token

	defer client.sessionToken.mx.Unlock()

	return token
}
