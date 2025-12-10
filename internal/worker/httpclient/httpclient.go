package httpclient

import (
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

type Config struct {
	UserAgent      string
	ProxyURL       string        // empty = no proxy
	RequestDelayMs int           // minimum delay between requests
	Timeout        time.Duration // request timeout
}

type Client struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	userAgent  string
}

func New(cfg Config) (*Client, error) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment, // default proxy from environment
	}

	// override proxy if specified in config
	if cfg.ProxyURL != "" {
		proxyURL, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transport,
	}

	var limiter *rate.Limiter
	if cfg.RequestDelayMs > 0 {
		delay := time.Duration(cfg.RequestDelayMs) * time.Millisecond
		limiter = rate.NewLimiter(rate.Every(delay), 1)
	}

	return &Client{
		httpClient: client,
		limiter:    limiter,
		userAgent:  cfg.UserAgent,
	}, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.limiter != nil {
		if err := c.limiter.Wait(req.Context()); err != nil {
			return nil, err
		}
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	return c.httpClient.Do(req)
}
