package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultRetries = 0
	defaultTimeout = 30 * time.Second
)

type Service interface {
	Do(ctx context.Context, method, endpoint string, body io.Reader, params Params) (*Response, error)
	Clone(options ...Option) (Service, error)
}

type client struct {
	Retries    int
	Token      string
	Timeout    time.Duration
	BaseURL    *url.URL
	httpClient *http.Client
}

type Option func(Service)

type Params map[string]string

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func New(baseURL string, options ...Option) (Service, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	c := &client{
		Retries: defaultRetries,
		Timeout: defaultTimeout,
		BaseURL: parsedURL,
	}

	for _, option := range options {
		option(c)
	}

	c.httpClient = &http.Client{Timeout: c.Timeout}

	return c, nil
}

func WithRetries(retries int) Option {
	return func(c Service) {
		if client, ok := c.(*client); ok && retries > 0 {
			client.Retries = retries
		}
	}
}

func WithToken(token string) Option {
	return func(c Service) {
		if client, ok := c.(*client); ok && token != "" {
			client.Token = token
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c Service) {
		if client, ok := c.(*client); ok && timeout > 0 {
			client.Timeout = timeout
		}
	}
}

func (c *client) Do(ctx context.Context, method, endpoint string, body io.Reader, params Params) (*Response, error) {
	req, err := c.newRequest(ctx, method, endpoint, body, params)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *client) newRequest(ctx context.Context, method, endpoint string, body io.Reader, params Params) (*http.Request, error) {
	u := c.BaseURL.ResolveReference(&url.URL{Path: endpoint})

	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create an http request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *client) do(req *http.Request) (*Response, error) {
	var lastErr error
	attempts := c.Retries + 1

	for i := range attempts {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d failed because of network error: %w", i+1, err)
			time.Sleep(time.Duration(i) * 100 * time.Microsecond)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= http.StatusInternalServerError {
			lastErr = fmt.Errorf("attempt %d failed because of server error %d", i+1, resp.StatusCode)
			time.Sleep(time.Duration(i) * 100 * time.Microsecond)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			time.Sleep(time.Duration(i) * 100 * time.Microsecond)
			continue
		}

		return &Response{
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Body:       body,
		}, nil
	}

	return nil, fmt.Errorf("all attempts failed after %d retries: %w", c.Retries, lastErr)
}

func (c *client) Clone(options ...Option) (Service, error) {
	clone := &client{
		Retries: c.Retries,
		Token:   c.Token,
		Timeout: c.Timeout,
		BaseURL: c.BaseURL,
	}

	for _, option := range options {
		option(clone)
	}

	clone.httpClient = &http.Client{Timeout: clone.Timeout}

	return clone, nil
}

func (r *Response) Unmarshal(v any) error {
	return json.Unmarshal(r.Body, v)
}
