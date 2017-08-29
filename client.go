package now

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
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

// ErrZeitResponse represents the body returned on error
type ErrZeitResponse struct {
	StatusCode int
	Err        *ZeitError `json:"err,omitempty"`
	AltErr     *ZeitError `json:"error,omitempty"`
}

// ZeitError returns the ZeitError depending on which response was given from the api
// Sometimes Zeit returns `err` sometimes `error`
// XXX: Reported to Zeit #now slack channel on Aug 24, 2017
func (e ErrZeitResponse) ZeitError() *ZeitError {
	if e.Err != nil {
		return e.Err
	}
	if e.AltErr != nil {
		return e.AltErr
	}
	return nil
}

// ZeitError contains the api-specific error response fields
type ZeitError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

// NewFileRequest performs an authenticated file upload for the given params
func (c Client) NewFileRequest(method, path string, file *os.File, v interface{}, headers *map[string]string) ClientError {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = c.URL + path

	req, err := http.NewRequest(method, path, file)
	if err != nil {
		return NewError(err.Error())
	}

	return c.performRequest(req, headers, v)
}

// NewRequest performs an authenticated request for the given params
func (c Client) NewRequest(method, path string, body interface{}, v interface{}, headers *map[string]string) ClientError {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = c.URL + path

	var req *http.Request
	var rErr error
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return NewError(err.Error())
		}
		req, rErr = http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
	} else {
		req, rErr = http.NewRequest(method, path, nil)
	}
	if rErr != nil {
		return NewError(rErr.Error())
	}

	return c.performRequest(req, headers, v)
}

func (c Client) performRequest(req *http.Request, headers *map[string]string, v interface{}) ClientError {
	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("User-Agent", "go-now")
	req.Header.Set("Authorization", "Bearer "+c.secret)

	// Optionally add teamID to every request
	if c.teamID != "" {
		q := req.URL.Query()
		q.Add("teamId", c.teamID)
		req.URL.RawQuery = q.Encode()
	}

	// Perform the request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return NewError(err.Error())
	}
	defer res.Body.Close()

	// Read and triage the response
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return NewError(err.Error())
	}
	switch res.StatusCode {
	case 200, 202, 201, 204:
		if v != nil && len(resBody) > 0 {
			err := json.Unmarshal(resBody, v)
			if err != nil {
				return NewError("Failed to read response")
			}
		}
		return nil
	case 304:
		return nil
	default:
		zeitErrResp := ErrZeitResponse{}
		if len(resBody) > 0 {
			marshalErr := json.Unmarshal(resBody, &zeitErrResp)
			if marshalErr != nil {
				return NewError("Invalid API response")
			}
		}
		return NewZeitError(res.StatusCode, zeitErrResp.ZeitError())
	}
}
