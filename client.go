package now

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// version contains the current package version
const version = "dev"

const apiURL = "https://api.zeit.co"

// Client contains all methods used for making API requests
type Client struct {
	secret     string
	URL        string
	HTTPClient *http.Client
}

// Authenticated returns whether the secret value is set
func (c Client) Authenticated() bool {
	return c.secret != ""
}

// SetHTTPClient overrides the default HTTP client used
func (c Client) SetHTTPClient(h *http.Client) {
	c.HTTPClient = h
}

// ErrResponse represents the body returned on error
type ErrResponse struct {
	StatusCode int
	Response   APIError `json:"err"`
}

func (e ErrResponse) Error() string {
	return fmt.Sprintf("%s (%d): %s", e.Response.Code, e.StatusCode, e.Response.Message)
}

// APIError contains the error response fields
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

// NewRequest performs an authenticated request for the given params
func (c Client) NewRequest(method, path string, body interface{}, v interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = c.URL + path

	var req *http.Request
	var err error
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req, err = http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "go-now@"+version)
	req.Header.Set("Authorization", "Bearer "+c.secret)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 204:
		if v != nil {
			return json.Unmarshal(resBody, v)
		}
		return nil
	default:
		apiErr := ErrResponse{StatusCode: res.StatusCode}
		marshalErr := json.Unmarshal(resBody, &apiErr)
		if marshalErr != nil {
			return marshalErr
		}
		return apiErr
	}
}
