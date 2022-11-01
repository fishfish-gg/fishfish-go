package fishfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (client *FishFishClient) fetchDomains() error {
	var err error
	var resp http.Response
	sessionToken := client.getSessionToken()

	if len(sessionToken) > 0 {
		httpResp, httpErr := client.authenticatedRequest("domains", "GET", "{}")

		resp = *httpResp
		err = httpErr
	} else {
		url := client.getAPIUrl("domains")
		httpResp, httpErr := http.Get(url)

		resp = *httpResp
		err = httpErr
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("API Returned non-200 status code.\nCode received: %d", resp.StatusCode))
		return err
	}

	var domains []string

	decoder := json.NewDecoder(resp.Body)

	if err = decoder.Decode(&domains); err != nil {
		return err
	}

	client.domainCache.mx.Lock()

	client.domainCache.entries = domains

	defer client.domainCache.mx.Unlock()

	return err
}

func (client *FishFishClient) GetDomains() []string {
	client.domainCache.mx.Lock()

	domains := make([]string, len(client.domainCache.entries))

	for index := range client.domainCache.entries {
		domains[index] = client.domainCache.entries[index]
	}

	defer client.domainCache.mx.Unlock()

	return domains
}
