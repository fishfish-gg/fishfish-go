package fishfish

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func DefaultConfig() Config {
	return Config{
		CacheInterval: 5000,
	}
}

func New(config Config) *FishFishClient {
	var client *FishFishClient
	client = &FishFishClient{
		config:     config,
		url:        fmt.Sprintf("https://api.fishfish.gg/v%d/", APIVersion),
		httpClient: &http.Client{},
	}

	client.fetchDomains()

	if len(config.Auth) > 0 {
		client.updateSessionToken()

		// Get session token ticker
		ctx2, cancel2 := context.WithCancel(context.Background())
		ticker2 := time.NewTicker(time.Hour * time.Duration(1))

		go func() {
			for {
				select {
				case <-ticker2.C:
					client.updateSessionToken()
				case <-ctx2.Done():
					return
				}
			}
		}()

		getNewSessionTicker := ClientTicker{
			ticker:    ticker2,
			ctxCancel: cancel2,
		}

		client.tickers.getNewSessionTicker = getNewSessionTicker
	}

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

	client.tickers.fetchDomainsTicker = fetchDomainsTicker

	return client
}

func (client *FishFishClient) Close() {
	FDT := client.tickers.fetchDomainsTicker
	GNST := client.tickers.getNewSessionTicker

	FDT.ticker.Stop()
	FDT.ctxCancel()

	GNST.ticker.Stop()
	GNST.ctxCancel()
}
