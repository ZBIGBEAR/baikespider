package http

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
)

type Config struct {
	RetryMax     int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
	Timeout      time.Duration
}

type HTTP struct {
	config Config
	client *retryablehttp.Client
}

func Default() *HTTP {
	return New(Config{
		RetryMax:     3,
		RetryWaitMax: 10 * time.Second,
		RetryWaitMin: 2 * time.Second,
		Timeout:      5 * time.Second,
	})
}

func New(config Config) *HTTP {
	client := retryablehttp.NewClient()
	client.Logger = nil
	client.RetryMax = config.RetryMax
	client.RetryWaitMin = config.RetryWaitMin
	client.RetryWaitMax = config.RetryWaitMax
	client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			return false, err
		}
		if resp.StatusCode >= 500 {
			return true, nil
		}
		return false, err
	}
	client.HTTPClient.Timeout = config.Timeout
	client.HTTPClient.Transport = apmhttp.WrapRoundTripper(client.HTTPClient.Transport)
	return &HTTP{config, client}
}

func (h *HTTP) Do(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	req, err := retryablehttp.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		apm.CaptureError(req.Context(), err).Send()
		return nil, err
	}
	return resp, nil
}

func (h *HTTP) DoWithAuthorization(ctx context.Context, authorization, method, url string, body []byte) (*http.Response, error) {
	req, err := retryablehttp.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", authorization)
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		apm.CaptureError(req.Context(), err).Send()
		return nil, err
	}
	return resp, nil
}
