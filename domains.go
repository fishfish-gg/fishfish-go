package fishfish

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type DomainCache struct {
	entries []string
	mx      sync.Mutex
}

func (client *Client) syncDomains() error {
	var err error
	sessionToken := client.getSessionToken()

	apiURL := client.getAPIUrl("domains")
	body := bytes.NewBufferString("{}")

	req, _ := http.NewRequest("GET", apiURL, body)

	if sessionToken != "" {
		req.Header.Set("Authorization", sessionToken)
	}

	resp, err := client.httpClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("API Returned non-200 status code.\nCode received: %d\nAPI Response: %s", resp.StatusCode, resp.Body))
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

func (client *Client) GetDomains() []string {
	client.domainCache.mx.Lock()

	domains := make([]string, len(client.domainCache.entries))

	domains = append(domains, client.domainCache.entries...)

	defer client.domainCache.mx.Unlock()

	return domains
}
