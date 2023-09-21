// Package chttp provide a simple http client
package chttp

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/net/http2"
)

var mDefault = New()

var Debug = false

type Client struct {
	client *retryablehttp.Client
}

const (
	defaultTimeout      = 60 * time.Second
	defaultRetryWaitMin = 100 * time.Millisecond // nolint
	defaultRetryWaitMax = 10 * time.Second
	defaultRetryMax     = 5
)

func Default() *Client {
	return mDefault
}

func getDefaultLogger() interface{} {
	if Debug {
		return log.Default()
	}
	return nil
}

func New(options ...ClientOption) *Client {
	transport := cleanhttp.DefaultPooledTransport()
	_ = http2.ConfigureTransport(transport)

	c := &Client{
		client: &retryablehttp.Client{
			HTTPClient: &http.Client{
				Timeout:   defaultTimeout,
				Transport: transport,
			},
			Logger:       getDefaultLogger(),
			RetryWaitMin: defaultRetryWaitMin,
			RetryWaitMax: defaultRetryWaitMax,
			RetryMax:     defaultRetryMax,
			CheckRetry:   retryablehttp.DefaultRetryPolicy,
			Backoff:      retryablehttp.DefaultBackoff,
			ErrorHandler: func(resp *http.Response, err error, numTries int) (*http.Response, error) {
				return resp, err
			},
		},
	}

	for _, o := range options {
		o.applyToClient(c)
	}

	return c
}

func (c *Client) do(r *internalRequest) *Response {
	res, err := c.client.Do(r)
	return makeResponse(res, err)
}

// Do same as http.Client Do
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	r, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}
	res := c.do(r)
	return res.Response, res.err
}

func (c *Client) Call(req *http.Request) *Response {
	r, err := retryablehttp.FromRequest(req)
	if err != nil {
		return makeResponse(nil, err)
	}
	return c.do(r)
}

func (c *Client) Post(ctx context.Context, reqURL string, data any, options ...RequestOption) *Response {
	reqBody, mime, err := marshalRequestData(data)
	if err != nil {
		return makeResponse(nil, err)
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, reqURL, reqBody)
	if err != nil {
		return makeResponse(nil, err)
	}

	if len(mime) > 0 {
		req.Header.Set("Content-Type", mime)
	}

	for _, o := range options {
		o.applyToRequest(req.Request)
	}

	return c.do(req)
}

func (c *Client) Get(ctx context.Context, reqURL string, query url.Values, options ...RequestOption) *Response {
	if len(query) > 0 {
		if strings.Contains(reqURL, "?") {
			reqURL += "&" + query.Encode()
		} else {
			reqURL += "?" + query.Encode()
		}
	}
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return makeResponse(nil, err)
	}
	for _, o := range options {
		o.applyToRequest(req.Request)
	}
	return c.do(req)
}

func (c *Client) PostForm(ctx context.Context, reqURL string, values url.Values, options ...RequestOption) *Response {
	return c.Post(ctx, reqURL, Form{Values: values}, options...)
}
