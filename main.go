package fishfish

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Auth          string
	CacheInterval int
}

type domainCache struct {
	entries []string
	mx      sync.Mutex
}

type sessionToken struct {
	token string
	mx    sync.Mutex
}

type ClientTicker struct {
	ticker    *time.Ticker
	ctxCancel context.CancelFunc
}

type ClientTickers struct {
	fetchDomainsTicker  ClientTicker
	getNewSessionTicker ClientTicker
}

type FishFishClient struct {
	config       Config
	sessionToken sessionToken
	url          string
	domainCache  domainCache
	tickers      ClientTickers
	httpClient   *http.Client
}

type TokenResponse struct {
	token   string
	expires int
}

// Target API version
var APIVersion = 1

func (client *FishFishClient) getAPIUrl(path string) string {
	return fmt.Sprintf("%s%s", client.url, path)
}

func DefaultConfig() Config {
	return Config{
		CacheInterval: 5000,
	}
}

func (client *FishFishClient) fetchDomains() (err error) {
	url := client.getAPIUrl("domains")
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("API Returned non-200 status code.\nCode received: %d", resp.StatusCode))
		return
	}

	var domains []string

	decoder := json.NewDecoder(resp.Body)

	if err = decoder.Decode(&domains); err != nil {
		return
	}

	client.domainCache.mx.Lock()

	client.domainCache.entries = domains

	defer client.domainCache.mx.Unlock()

	return
}

func (client *FishFishClient) fetchSessionToken() (err error, token string) {
	req, _ := http.NewRequest("GET", client.getAPIUrl("users/@me/tokens"), nil)
	req.Header.Set("Authorization", client.config.Auth)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	tokenResponse := TokenResponse{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return
	}

	token = tokenResponse.token

	return
}

func (client *FishFishClient) getSessionToken() (token string) {
	client.sessionToken.mx.Lock()

	token = client.sessionToken.token

	defer client.sessionToken.mx.Unlock()

	return
}

func (client *FishFishClient) authenticatedRequest(path, requestType string) (*http.Response, error) {
	req, _ := http.NewRequest(strings.ToUpper(requestType), client.getAPIUrl(path), nil)
	req.Header.Set("Authorization", client.getSessionToken())
	resp, err := client.httpClient.Do(req)

	return resp, err
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

func (client *FishFishClient) Kill() {
	FDT := client.tickers.fetchDomainsTicker
	GNST := client.tickers.getNewSessionTicker

	FDT.ticker.Stop()
	FDT.ctxCancel()

	GNST.ticker.Stop()
	GNST.ctxCancel()
}

func New(config Config) *FishFishClient {
	var client *FishFishClient
	client = &FishFishClient{
		config:     config,
		url:        fmt.Sprintf("https://api.fishfish.gg/v%d/", APIVersion),
		httpClient: &http.Client{},
	}

	client.fetchDomains()

	// Fetch domains ticker
	ctx1, cancel1 := context.WithCancel(context.Background())

	ticker1 := time.NewTicker(time.Millisecond * time.Duration(config.CacheInterval))

	go func() {
		for {
			select {
			case <-ticker1.C:
				client.fetchDomains()
			case <-ctx1.Done():
				return
			}
		}
	}()
	fetchDomainsTicker := ClientTicker{
		ticker:    ticker1,
		ctxCancel: cancel1,
	}

	// Get session token ticker
	ctx2, cancel2 := context.WithCancel(context.Background())

	ticker2 := time.NewTicker(time.Hour * time.Duration(1))

	go func() {
		for {
			select {
			case <-ticker2.C:
				err, sessionToken := client.fetchSessionToken()
				if err != nil {
					log.Fatal(err)
				}

				client.sessionToken.mx.Lock()

				client.sessionToken.token = sessionToken

				defer client.sessionToken.mx.Unlock()
			case <-ctx2.Done():
				return
			}
		}
	}()
	getNewSessionTicker := ClientTicker{
		ticker:    ticker2,
		ctxCancel: cancel2,
	}

	// Add tickers to client
	client.tickers.fetchDomainsTicker = fetchDomainsTicker
	client.tickers.getNewSessionTicker = getNewSessionTicker

	return client
}
