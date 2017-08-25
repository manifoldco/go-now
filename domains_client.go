package now

import (
	"fmt"
)

const domainsEndpoint = "/domains"

// DomainsClient contains the methods for the Domain API
type DomainsClient struct {
	client *Client
}

// New creates a new Domain
func (c DomainsClient) New(domainName string, external bool) (Domain, ClientError) {
	return c.NewFromParams(DomainParams{
		Name:       domainName,
		IsExternal: external,
	})
}

// NewFromParams creates a new Domain from params
func (c DomainsClient) NewFromParams(params DomainParams) (Domain, ClientError) {
	d := Domain{}
	err := c.client.NewRequest("POST", domainsEndpoint, params, &d, nil)
	return d, err
}

// DomainParams contains all fields for domain create
type DomainParams struct {
	Name       string `json:"name"`
	IsExternal bool   `json:"isExternal"`
}

// List retrieves a list of all the domains under the account
func (c DomainsClient) List() ([]Domain, ClientError) {
	d := &domainListResponse{}
	err := c.client.NewRequest("GET", domainsEndpoint, nil, d, nil)
	return d.Domains, err
}

type domainListResponse struct {
	Domains []Domain `json:"domains"`
}

// Delete deletes the domain by its ID
func (c DomainsClient) Delete(domainName string) ClientError {
	return c.client.NewRequest("DELETE", fmt.Sprintf("%s/%s", domainsEndpoint, domainName), nil, nil, nil)
}
