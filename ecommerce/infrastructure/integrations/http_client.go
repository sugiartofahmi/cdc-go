package integrations

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpClient struct {
	client  *http.Client
	headers map[string]string
	params  map[string]string
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: &http.Client{Timeout: 30 * time.Second},
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		params: make(map[string]string),
	}
}

func (h *HttpClient) SetTimeout(timeout time.Duration) {
	h.client.Timeout = timeout
}

func (h *HttpClient) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		h.headers[k] = v
	}
}

func (h *HttpClient) SetParams(params map[string]string) {
	for k, v := range params {
		h.params[k] = v
	}
}

func (h *HttpClient) Get(rawURL string) ([]byte, error) {
	fullURL, err := h.buildURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	h.setRequestHeaders(req)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return respBody, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (h *HttpClient) Post(rawURL string, body io.Reader) ([]byte, error) {
	fullURL, err := h.buildURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequest("POST", fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	h.setRequestHeaders(req)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return respBody, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (h *HttpClient) buildURL(rawURL string) (string, error) {
	if len(h.params) == 0 {
		return rawURL, nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := parsed.Query()
	for k, v := range h.params {
		q.Set(k, v)
	}
	parsed.RawQuery = q.Encode()

	return parsed.String(), nil
}

func (h *HttpClient) setRequestHeaders(req *http.Request) {
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}
}
