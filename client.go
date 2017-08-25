package now

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const apiURL = "https://api.zeit.co"

// Client contains all methods used for making API requests
type Client struct {
	secret     string
	teamID     string
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

// ErrAPIResponse represents the body returned on error
type ErrAPIResponse struct {
	StatusCode int
	Err        *APIError `json:"err,omitempty"`
	AltErr     *APIError `json:"error,omitempty"`
}

// APIError returns the APIError depending on which response was given from the api
// Sometimes Zeit returns `err` sometimes `error`
// XXX: Reported to Zeit #now slack channel on Aug 24, 2017
func (e ErrAPIResponse) APIError() APIError {
	err := e.AltErr
	if e.Err != nil {
		err = e.Err
	}
	return *err
}

// ErrResponse represents the body returned on error
type ErrResponse struct {
	StatusCode int      `json:"status_code"`
	APIError   APIError `json:"err"`
}

// Error returns the error string
func (e ErrResponse) Error() string {
	err := e.APIError
	return fmt.Sprintf("%s (%d): %s", err.Code, e.StatusCode, err.Message)
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
	var rErr error
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req, rErr = http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
	} else {
		req, rErr = http.NewRequest(method, path, nil)
	}
	if rErr != nil {
		return rErr
	}

	req.Header.Set("User-Agent", "go-now")
	req.Header.Set("Authorization", "Bearer "+c.secret)
	req.Header.Set("Content-Type", "application/json")

	// Optionally add teamID to every request
	if c.teamID != "" {
		q := req.URL.Query()
		q.Add("teamId", c.teamID)
		req.URL.RawQuery = q.Encode()
	}

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
		apiErr := ErrAPIResponse{}
		marshalErr := json.Unmarshal(resBody, &apiErr)
		if marshalErr != nil {
			return marshalErr
		}
		return ErrResponse{
			StatusCode: res.StatusCode,
			APIError:   apiErr.APIError(),
		}
	}
}
