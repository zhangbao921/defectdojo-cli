package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	BaseURL  string
	APIKey   string
	HTTP     *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return &Client{
		BaseURL: baseURL + "api/v2/",
		APIKey:  apiKey,
		HTTP:    &http.Client{},
	}
}

func (c *Client) doRequest(method, path string, body io.Reader, contentType string, params url.Values) ([]byte, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+c.APIKey)
	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

func (c *Client) Get(path string, params url.Values) ([]byte, error) {
	return c.doRequest(http.MethodGet, path, nil, "", params)
}

type PaginatedResult struct {
	Count   int                `json:"count"`
	Next    *string            `json:"next"`
	Results []json.RawMessage  `json:"results"`
}

func (c *Client) ListAll(path string, params url.Values) ([]json.RawMessage, error) {
	if params == nil {
		params = url.Values{}
	}
	if params.Get("limit") == "" {
		params.Set("limit", "200")
	}

	var all []json.RawMessage
	page := 1

	for {
		p := url.Values{}
		for k, v := range params {
			p[k] = v
		}
		p.Set("page", fmt.Sprintf("%d", page))

		data, err := c.Get(path, p)
		if err != nil {
			return nil, err
		}

		var pr PaginatedResult
		if err := json.Unmarshal(data, &pr); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		all = append(all, pr.Results...)
		if pr.Next == nil || *pr.Next == "" {
			break
		}
		page++
	}
	return all, nil
}

func (c *Client) GetByID(path, id string) ([]byte, error) {
	return c.Get(strings.TrimRight(path, "/")+"/"+id+"/", nil)
}

func (c *Client) Post(path string, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode: %w", err)
	}
	// Ensure trailing slash for Django's APPEND_SLASH (prevents POST -> GET redirect)
	path = strings.TrimRight(path, "/") + "/"
	return c.doRequest(http.MethodPost, path, &buf, "application/json", nil)
}

func (c *Client) Put(path, id string, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode: %w", err)
	}
	return c.doRequest(http.MethodPut, strings.TrimRight(path, "/")+"/"+id+"/", &buf, "application/json", nil)
}

func (c *Client) Patch(path, id string, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode: %w", err)
	}
	return c.doRequest(http.MethodPatch, strings.TrimRight(path, "/")+"/"+id+"/", &buf, "application/json", nil)
}

func (c *Client) Delete(path, id string) error {
	_, err := c.doRequest(http.MethodDelete, strings.TrimRight(path, "/")+"/"+id+"/", nil, "", nil)
	return err
}

func (c *Client) PostMultipart(path string, fields map[string]string, fileField, filePath string) ([]byte, error) {
	// Ensure trailing slash for Django
	path = strings.TrimRight(path, "/") + "/"
	var buf bytes.Buffer
	mp := multipart.NewWriter(&buf)

	for k, v := range fields {
		mp.WriteField(k, v)
	}

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		defer file.Close()

		w, err := mp.CreateFormFile(fileField, filepath.Base(filePath))
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		if _, err := io.Copy(w, file); err != nil {
			return nil, fmt.Errorf("failed to copy file: %w", err)
		}
	}

	contentType := mp.FormDataContentType()
	mp.Close()

	return c.doRequest(http.MethodPost, path, &buf, contentType, nil)
}
