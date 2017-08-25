package now

import (
	"fmt"
)

// DomainsClient contains the methods for the Domain API
type DomainsClient struct {
	client *Client
}

// New creates a new Domain
func (c DomainsClient) New(domainName string, external bool) (Domain, error) {
	d := Domain{}
	params := domainParams{
		Name:       domainName,
		IsExternal: external,
	}
	err := c.client.NewRequest("POST", "/domains", params, &d)
	return d, err
}

type domainParams struct {
	Name       string `json:"name"`
	IsExternal bool   `json:"isExternal"`
}

// List retrieves a list of all the domains under the account
func (c DomainsClient) List() ([]Domain, error) {
	d := &domainListResponse{}
	err := c.client.NewRequest("GET", "/domains", nil, d)
	return d.Domains, err
}

type domainListResponse struct {
	Domains []Domain `json:"domains"`
}

// Delete deletes the domain by its ID
func (c DomainsClient) Delete(domainName string) error {
	return c.client.NewRequest("DELETE", fmt.Sprintf("/domains/%s", domainName), nil, nil)
}
