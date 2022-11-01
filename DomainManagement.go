package fishfish

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (client *FishFishClient) AddDomain(domain string, options CreateDomainBody) error {
	sessionToken := client.getSessionToken()
	if len(sessionToken) <= 0 {
		return errors.New("This function requires authentication!")
	}

	jsonBody, err := json.Marshal(options)
	if err != nil {
		return err
	}

	_, err = client.authenticatedRequest(fmt.Sprintf("domains/%s", domain), "POST", string(jsonBody))

	if err != nil {
		client.fetchDomains()
	}
	return err
}

func (client *FishFishClient) UpdateDomain(domain string, options UpdateDomainBody) error {
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

	_, err = client.authenticatedRequest(fmt.Sprintf("domains/%s", domain), "PATCH", string(jsonBody))

	return err
}

func (client *FishFishClient) DeleteDomain(domain string) error {
	sessionToken := client.getSessionToken()

	if len(sessionToken) <= 0 {
		return errors.New("This function requires authentication!")
	}

	_, err := client.authenticatedRequest(fmt.Sprintf("domains/%s", domain), "DELETE", "{}")
	return err
}
