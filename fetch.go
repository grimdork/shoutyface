package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Request an endpoint without a response body.
func Request(method, path string, headers map[string]string) (*http.Response, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	url := cfg.ServerURL + "/api/" + path
	headers["token"] = cfg.AdminToken
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	c := &http.Client{Timeout: time.Second * 10}
	res, err := c.Do(req)
	if err != nil {
		return res, err
	}

	return res, res.Body.Close()
}

// RequestJSON requests an endpoint and converts a JSON response.
// The response body will be closed.
func RequestJSON(method, path string, headers map[string]string, out interface{}) (*http.Response, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	url := cfg.ServerURL + "/api/" + path
	headers["token"] = cfg.AdminToken
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	c := &http.Client{Timeout: time.Second * 10}
	res, err := c.Do(req)
	if err != nil {
		return res, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusCreated {
		return res, nil
	}

	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return res, json.Unmarshal(buf, out)
}
