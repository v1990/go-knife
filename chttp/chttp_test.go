package chttp

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	Debug = true
}

func TestClient_Get(t *testing.T) {
	server := newTestServer()
	baseURL := server.URL

	c := New(
		WithHTTPTimeout(time.Second*60),
		WithRetryLimit(5, time.Millisecond*10, time.Second),
		// WithRetryPolicy(retryablehttp.ErrorPropagatedRetryPolicy),
		WithDialer(net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
			// DualStack: true,
		}),
	)
	t.Run("GET", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		query := url.Values{}
		query.Set("k2", "v2")
		resp := c.Get(ctx, baseURL+"/echo?k1=v1", query)
		err := resp.CheckStatus(200).Err()
		require.NoError(t, err)
	})
	t.Run("POST", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		in := map[string]string{"a": "b"}
		resp := c.Post(ctx, baseURL+"/echo", in)

		out := make(map[string]string)
		err := resp.CheckStatus(200).UnmarshalBodyJSON(&out)
		require.NoError(t, err)
	})
	t.Run("PostForm", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		form := make(url.Values)
		form.Add("k", "v")
		resp := c.PostForm(ctx, baseURL+"/echo", form)
		err := resp.CheckStatus(200).Err()
		require.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp := c.Get(ctx, baseURL+"/error", nil)
		err := resp.CheckStatus(200).Err()
		t.Logf("err(%T): %+v", err, err)
		require.Error(t, err)
	})
	t.Run("bad addr", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		resp := c.Get(ctx, "http://10.0.1.1:22", nil)
		err := resp.Err()
		t.Logf("err(%T): %+v", err, err)
	})
}

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.Copy(w, r.Body)
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	})

	return httptest.NewServer(mux)
}
