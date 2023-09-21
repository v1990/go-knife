package chttp

import (
	"log"
	"net"
	"net/http"
	"time"
)

// ClientOption Options for setting the Client
type ClientOption interface {
	applyToClient(c *Client)
}

type clientOptionFunc func(c *Client)

func (f clientOptionFunc) applyToClient(c *Client) { f(c) }

func WithHTTPTimeout(to time.Duration) ClientOption {
	return clientOptionFunc(func(c *Client) {
		c.client.HTTPClient.Timeout = to
	})
}

func WithLogger(logger Logger) ClientOption {
	return clientOptionFunc(func(c *Client) {
		c.client.Logger = logger
	})
}

func WithLeveledLogger(logger LeveledLogger) ClientOption {
	return clientOptionFunc(func(c *Client) {
		c.client.Logger = logger
	})
}

func WithDefaultStdLogger() ClientOption {
	return WithLogger(log.Default())
}

func WithRetryLimit(retryMax int, retryWaitMin, retryWaitMax time.Duration) ClientOption {
	return clientOptionFunc(func(c *Client) {
		c.client.RetryMax = retryMax
		c.client.RetryWaitMin = retryWaitMin
		c.client.RetryWaitMax = retryWaitMax
	})
}

// WithRetryPolicy sets the retry policy for the client.
func WithRetryPolicy(policy CheckRetry) ClientOption {
	return clientOptionFunc(func(c *Client) {
		c.client.CheckRetry = policy
	})
}

func WithBackoff(backoff Backoff) ClientOption {
	return clientOptionFunc(func(c *Client) {
		c.client.Backoff = backoff
	})
}

func WithErrorHandler(handler ErrorHandler) ClientOption {
	return clientOptionFunc(func(c *Client) {
		c.client.ErrorHandler = handler
	})
}

func WithDialer(dialer net.Dialer) ClientOption {
	return clientOptionFunc(func(c *Client) {
		if transport, ok := c.client.HTTPClient.Transport.(*http.Transport); ok {
			transport.DialContext = dialer.DialContext
		} else {
			panic("not supported")
		}
	})
}

// RequestOption Options for setting the http.Request
type RequestOption interface {
	applyToRequest(r *http.Request)
}

type requestOptionFunc func(r *http.Request)

func (f requestOptionFunc) applyToRequest(r *http.Request) { f(r) }

func WithHeader(kvs ...string) RequestOption {
	return requestOptionFunc(func(r *http.Request) {
		if len(kvs)%2 != 0 {
			panic("len(kvs) must be even")
		}
		for len(kvs) >= 2 {
			k := kvs[0]
			v := kvs[1]
			kvs = kvs[2:]
			r.Header.Set(k, v)
		}
	})
}

func WithContentType(mime string) RequestOption {
	return WithHeader("Content-Type", mime)
}
