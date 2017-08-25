package now

import (
	"net/http"
	"time"
)

const defaultHTTPTimeout = 80 * time.Second

// Now contains all the methods required for interacting with Zeit Now's API
type Now struct {
	client      *Client
	Deployments *DeploymentsClient

// SetTeamID updates the client's global team_id value
func (n Now) SetTeamID(teamID string) {
	n.client.teamID = teamID
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
	n.Deployments = &DeploymentsClient{client: n.client}
	return &n
}
