package now

import (
	"net/http"
	"time"
)

const defaultHTTPTimeout = 80 * time.Second

// Now contains all the methods required for interacting with Zeit Now's API
type Now struct {
	client     *Client
	Deployment *DeploymentClient
}

// New returns an authenticated Now api client
func New(secret string) *Now {
	n := Now{
		client: &Client{
			secret: secret,
			URL:    apiURL,
			HTTPClient: &http.Client{
				Timeout: defaultHTTPTimeout,
			},
		},
	}
	n.Deployment = &DeploymentClient{client: n.client}
	return &n
}
