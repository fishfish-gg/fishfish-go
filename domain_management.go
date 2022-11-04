package fishfish

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type CreateDomainBody struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Target      string `json:"target,omitempty"`
}

type UpdateDomainBody struct {
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Target      string `json:"target,omitempty"`
}

func (client *Client) AddDomain(domain string, options CreateDomainBody) error {
	sessionToken := client.getSessionToken()
	if len(sessionToken) <= 0 {
		return errors.New("This function requires authentication!")
	}

	jsonBody, err := json.Marshal(options)
	if err != nil {
		return err
	}

	apiURL := client.getAPIUrl(fmt.Sprintf("domains/%s", domain))
	body := bytes.NewBufferString(string(jsonBody))

	req, _ := http.NewRequest("PATCH", apiURL, body)
	req.Header.Set("Authorization", sessionToken)

	_, err = client.httpClient.Do(req)

	if err == nil {
		client.syncDomains()
	}

	return err
}

func (client *Client) UpdateDomain(domain string, options UpdateDomainBody) error {
	sessionToken := client.getSessionToken()

	if len(sessionToken) <= 0 {
		return errors.New("This function requires authentication!")
	}

	if !validCategory(options.Category) {
		return errors.New("Invalid category!")
	}

	jsonBody, err := json.Marshal(options)
	if err != nil {
		return err
	}

	apiURL := client.getAPIUrl(fmt.Sprintf("domains/%s", domain))
	body := bytes.NewBufferString(string(jsonBody))

	req, _ := http.NewRequest("PATCH", apiURL, body)
	req.Header.Set("Authorization", sessionToken)

	_, err = client.httpClient.Do(req)

	return err
}

func (client *Client) DeleteDomain(domain string) error {
	sessionToken := client.getSessionToken()

	if len(sessionToken) <= 0 {
		return errors.New("This function requires authentication!")
	}

	apiURL := client.getAPIUrl(fmt.Sprintf("domains/%s", domain))
	body := bytes.NewBufferString("{}")

	req, _ := http.NewRequest("DELETE", apiURL, body)
	req.Header.Set("Authorization", sessionToken)

	_, err := client.httpClient.Do(req)

	return err
}
