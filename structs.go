package fishfish

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	Auth          string
	CacheInterval int
	Permissions   []string
}

type DomainCache struct {
	entries []string
	mx      sync.Mutex
}

type SessionToken struct {
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
	sessionToken SessionToken
	url          string
	domainCache  DomainCache
	tickers      ClientTickers
	httpClient   *http.Client
}

type createTokenRequest struct {
	Permissions []string `json:"Permissions"`
}

type TokenResponse struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}

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
