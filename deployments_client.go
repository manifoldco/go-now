package now

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	endpointSync             = "/now/sync"
	endpointCreateDeployment = "/now/create"
	endpointDeployments      = "/now/deployments"
	endpointDeploymentsID    = "/now/deployments/%s"
)

// DeploymentsClient contains the methods for the Deployment API
type DeploymentsClient struct {
	client *Client
}

// New creates a new Deployment
func (c DeploymentsClient) New(params DeploymentParams) (IncompleteDeployment, ClientError) {
	d := IncompleteDeployment{}
	err := c.client.NewRequest("POST", endpointCreateDeployment, params, &d, nil)
	// TODO warn about invalid files, or size issues
	return d, err
}

// DeploymentParams contains all fields necessary to create a deployment
type DeploymentParams struct {
	Env         map[string]string `json:"env"`
	Public      bool              `json:"public"`
	ForceNew    bool              `json:"forceNew"`
	ForceSync   bool              `json:"forceSync"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"deploymentType"`
	Files       []FileInfo        `json:"files"`
}

// Upload performs an upload of the given file to the specified deployment
func (c DeploymentsClient) Upload(deploymentID, sha string, names []string, size int64, data *os.File) ClientError {
	headers := map[string]string{
		"Content-Type":        "application/octet-stream",
		"x-now-deployment-id": deploymentID,
		"x-now-sha":           sha,
		"x-now-file":          strings.Join(names, ","),
		"x-now-size":          strconv.Itoa(int(size)),
	}
	return c.client.NewFileRequest("POST", endpointSync, data, nil, &headers)
}

// Get retrieves a deployment by its ID
func (c DeploymentsClient) Get(ID string) (Deployment, ClientError) {
	d := Deployment{}
	err := c.client.NewRequest("GET", fmt.Sprintf(endpointDeploymentsID, ID), nil, &d, nil)
	return d, err
}

// Scale sets the scale of a deployment to the number provided
func (c DeploymentsClient) Scale(ID string, min, max int) (Deployment, ClientError) {
	d := Deployment{}
	err := c.client.NewRequest("POST", fmt.Sprintf(endpointDeploymentsID+"/instances", ID), ScaleParams{
		Min: min,
		Max: max,
	}, &d, nil)
	return d, err
}

// ScaleParams contains all fields necessary to set the scale of a deployment
type ScaleParams struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// Alias applies the supplied alias to the given deployment ID
func (c DeploymentsClient) Alias(ID, alias string) (Alias, ClientError) {
	a := Alias{Alias: alias}
	err := c.client.NewRequest("POST", fmt.Sprintf(endpointDeploymentsID+"/aliases", ID), DeploymentAliasParams{Alias: alias}, &a, nil)
	return a, err
}

// DeploymentAliasParams contains all fields for aliasing
type DeploymentAliasParams struct {
	Alias string `json:"alias"`
}

// ListAliases retrieves aliases of a deployment by its ID
func (c DeploymentsClient) ListAliases(ID string) ([]Alias, ClientError) {
	a := &deploymentListAliasResponse{}
	err := c.client.NewRequest("GET", fmt.Sprintf(endpointDeploymentsID+"/aliases", ID), nil, a, nil)
	return a.Aliases, err
}

type deploymentListAliasResponse struct {
	Aliases []Alias `json:"aliases"`
}

// Files retrieves files of a deployment by its ID
func (c DeploymentsClient) Files(ID string) ([]DeploymentContent, ClientError) {
	var contents []DeploymentContent
	var resp []json.RawMessage
	err := c.client.NewRequest("GET", fmt.Sprintf(endpointDeploymentsID+"/files", ID), nil, &resp, nil)
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
	err := c.client.NewRequest("GET", endpointDeployments, nil, d, nil)
	return d.Deployments, err
}

type deploymentListResponse struct {
	Deployments []Deployment `json:"deployments"`
}

// Delete deletes the deployment by its ID
func (c DeploymentsClient) Delete(ID string) ClientError {
	return c.client.NewRequest("DELETE", fmt.Sprintf(endpointDeploymentsID, ID), nil, nil, nil)
}
