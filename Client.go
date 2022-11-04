package fishfish

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Auth          string
	CacheInterval int
	Permissions   []string
}

type ClientTicker struct {
	ticker    *time.Ticker
	ctxCancel context.CancelFunc
}

type ClientTickers struct {
	syncDomainsTicker   ClientTicker
	getNewSessionTicker ClientTicker
}

type Client struct {
	config       Config
	sessionToken SessionToken
	url          string
	domainCache  DomainCache
	tickers      ClientTickers
	httpClient   *http.Client
}

var APIVersion = 1

func DefaultConfig() Config {
	return Config{
		CacheInterval: 5000,
	}
}

func New(config Config) *Client {
	var client *Client
	client = &Client{
		config:     config,
		url:        fmt.Sprintf("https://api.fishfish.gg/v%d/", APIVersion),
		httpClient: &http.Client{},
	}

	client.syncDomains()

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

	// Sync domains ticker
	ctx1, cancel1 := context.WithCancel(context.Background())
	ticker1 := time.NewTicker(time.Millisecond * time.Duration(config.CacheInterval))

	go func() {
		for {
			select {
			case <-ticker1.C:
				err := client.syncDomains()
				if err != nil {
					fmt.Println(err)
				}
			case <-ctx1.Done():
				return
			}
		}
	}()

	syncDomainsTicker := ClientTicker{
		ticker:    ticker1,
		ctxCancel: cancel1,
	}

	client.tickers.syncDomainsTicker = syncDomainsTicker

	return client
}

func (client *Client) Close() {
	FDT := client.tickers.syncDomainsTicker
	GNST := client.tickers.getNewSessionTicker

	FDT.ticker.Stop()
	FDT.ctxCancel()

	GNST.ticker.Stop()
	GNST.ctxCancel()
}
