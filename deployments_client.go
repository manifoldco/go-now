package now

import (
	"encoding/json"
	"errors"
	"fmt"
)

// DeploymentsClient contains the methods for the Deployment API
type DeploymentsClient struct {
	client *Client
}

// New creates a new Deployment
func (c DeploymentsClient) New(params map[string]interface{}) (Deployment, error) {
	d := Deployment{}
	err := c.client.NewRequest("POST", "/now/deployments", params, &d)
	return d, err
}

// Get retrieves a deployment by its ID
func (c DeploymentsClient) Get(ID string) (Deployment, error) {
	d := Deployment{}
	err := c.client.NewRequest("GET", fmt.Sprintf("/now/deployments/%s", ID), nil, &d)
	return d, err
}

// Alias applies the supplied alias to the given deployment ID
func (c DeploymentsClient) Alias(ID, alias string) (Alias, error) {
	a := Alias{}
	err := c.client.NewRequest("POST", fmt.Sprintf("/now/deployments/%s/aliases", ID), deploymentAliasParams{Alias: alias}, &a)
	return a, err
}

type deploymentAliasParams struct {
	Alias string `json:"alias"`
}

// ListAliases retrieves aliases of a deployment by its ID
func (c DeploymentsClient) ListAliases(ID string) ([]Alias, error) {
	a := &deploymentListAliasResponse{}
	err := c.client.NewRequest("GET", fmt.Sprintf("/now/deployments/%s/aliases", ID), nil, a)
	return a.Aliases, err
}

type deploymentListAliasResponse struct {
	Aliases []Alias `json:"aliases"`
}

// Files retrieves files of a deployment by its ID
func (c DeploymentsClient) Files(ID string) ([]DeploymentContent, error) {
	var contents []DeploymentContent
	var resp []json.RawMessage
	err := c.client.NewRequest("GET", fmt.Sprintf("/now/deployments/%s/files", ID), nil, &resp)
	for _, r := range resp {
		var obj map[string]interface{}

		// Extract the type field
		err := json.Unmarshal(r, &obj)
		if err != nil {
			return contents, err
		}

		// Unmarshal into appropriate type
		var content DeploymentContent
		switch obj["type"].(string) {
		case "directory":
			content = &DeploymentDir{}
		case "file":
			content = &DeploymentFile{}
		default:
			err = errors.New("Unknown file type")
			return contents, err
		}
		err = json.Unmarshal(r, &content)
		if err != nil {
			return contents, err
		}
		contents = append(contents, content)
	}
	return contents, err
}

// List retrieves a list of all the deployments under the account
func (c DeploymentsClient) List() ([]Deployment, error) {
	d := &deploymentListResponse{}
	err := c.client.NewRequest("GET", "/now/deployments", nil, d)
	return d.Deployments, err
}

type deploymentListResponse struct {
	Deployments []Deployment `json:"deployments"`
}

// Delete deletes the deployment by its ID
func (c DeploymentsClient) Delete(ID string) error {
	return c.client.NewRequest("DELETE", fmt.Sprintf("/now/deployments/%s", ID), nil, nil)
}
