package utils

import (
	"encoding/base64"
	"io"
	"net/http"
)

// HTTPRequestParams holds the parameters for the HTTP request
type HTTPRequestParams struct {
	Method   string
	URL      string
	Username string
	Password string
	Headers  map[string]string
	Body     io.Reader
}

// MakeRequest makes an HTTP request based on the given parameters
func MakeRequest(params HTTPRequestParams) (*http.Response, error) {
	client := &http.Client{}
	// Create the request
	req, err := http.NewRequest(params.Method, params.URL, params.Body)
	if err != nil {
		return nil, err
	}
	// Add basic auth if provided
	if params.Username != "" && params.Password != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(params.Username + ":" + params.Password))
		req.Header.Add("Authorization", "Basic "+auth)
	}
	// Add headers if provided
	for key, value := range params.Headers {
		req.Header.Add(key, value)
	}
	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
