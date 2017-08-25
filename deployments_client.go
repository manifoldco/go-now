package now

import (
	"encoding/json"
	"fmt"
)

const createDeploymentEndpoint = "/now/create"
const deploymentsEndpoint = "/now/deployments"

// DeploymentsClient contains the methods for the Deployment API
type DeploymentsClient struct {
	client *Client
}

// New creates a new Deployment
func (c DeploymentsClient) New(params map[string]interface{}) (Deployment, ClientError) {
	d := Deployment{}
	err := c.client.NewRequest("POST", createDeploymentEndpoint, params, &d, nil)
	return d, err
}

// Get retrieves a deployment by its ID
func (c DeploymentsClient) Get(ID string) (Deployment, ClientError) {
	d := Deployment{}
	err := c.client.NewRequest("GET", fmt.Sprintf("%s/%s", deploymentsEndpoint, ID), nil, &d, nil)
	return d, err
}

// Alias applies the supplied alias to the given deployment ID
func (c DeploymentsClient) Alias(ID, alias string) (Alias, ClientError) {
	a := Alias{Alias: alias}
	err := c.client.NewRequest("POST", fmt.Sprintf("%s/%s/aliases", deploymentsEndpoint, ID), DeploymentAliasParams{Alias: alias}, &a, nil)
	return a, err
}

// DeploymentAliasParams contains all fields for aliasing
type DeploymentAliasParams struct {
	Alias string `json:"alias"`
}

// ListAliases retrieves aliases of a deployment by its ID
func (c DeploymentsClient) ListAliases(ID string) ([]Alias, ClientError) {
	a := &deploymentListAliasResponse{}
	err := c.client.NewRequest("GET", fmt.Sprintf("%s/%s/aliases", deploymentsEndpoint, ID), nil, a, nil)
	return a.Aliases, err
}

type deploymentListAliasResponse struct {
	Aliases []Alias `json:"aliases"`
}

// Files retrieves files of a deployment by its ID
func (c DeploymentsClient) Files(ID string) ([]DeploymentContent, ClientError) {
	var contents []DeploymentContent
	var resp []json.RawMessage
	err := c.client.NewRequest("GET", fmt.Sprintf("%s/%s/files", deploymentsEndpoint, ID), nil, &resp, nil)
	for _, r := range resp {
		var obj map[string]interface{}

		// Extract the type field
		err := json.Unmarshal(r, &obj)
		if err != nil {
			return contents, NewError(err.Error())
		}

		// Unmarshal into appropriate type
		var content DeploymentContent
		switch obj["type"].(string) {
		case "directory":
			content = &DeploymentDir{}
		case "file":
			content = &DeploymentFile{}
		default:
			return contents, NewError("Unknown file type")
		}
		err = json.Unmarshal(r, &content)
		if err != nil {
			return contents, NewError(err.Error())
		}
		contents = append(contents, content)
	}
	return contents, err
}

// List retrieves a list of all the deployments under the account
func (c DeploymentsClient) List() ([]Deployment, ClientError) {
	d := &deploymentListResponse{}
	err := c.client.NewRequest("GET", deploymentsEndpoint, nil, d, nil)
	return d.Deployments, err
}

type deploymentListResponse struct {
	Deployments []Deployment `json:"deployments"`
}

// Delete deletes the deployment by its ID
func (c DeploymentsClient) Delete(ID string) ClientError {
	return c.client.NewRequest("DELETE", fmt.Sprintf("%s/%s", deploymentsEndpoint, ID), nil, nil, nil)
}
