package provider

// Source: https://github.com/BetterStackHQ/terraform-provider-better-uptime/blob/master/internal/provider/client.go

import (
	"context"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"log"
	"net/http"
	"strings"
)

type client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	userAgent  string
}

type option func(c *client)

func withHTTPClient(httpClient *http.Client) option {
	return func(c *client) {
		c.httpClient = httpClient
	}
}

func withUserAgent(userAgent string) option {
	return func(c *client) {
		c.userAgent = userAgent
	}
}

func newClient(baseURL, token string, opts ...option) (*client, error) {
	baseURL = strings.TrimSuffix(baseURL, "/")

	c := client{
		baseURL:    baseURL,
		token:      token,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c, nil
}

func (c *client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, path, nil)
}

func (c *client) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, path, body)
}

func (c *client) Patch(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	log.Print("XXX PATCH XXX")
	log.Print(path)
	return c.do(ctx, http.MethodPatch, path, body)
}

func (c *client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.do(ctx, http.MethodDelete, path, nil)
}

func (c *client) do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	if method == http.MethodPost || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}

	// TODO check for return code and return error if not in range 200-299
	return ctxhttp.Do(ctx, c.httpClient, req)
}
