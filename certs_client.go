package now

import (
	"fmt"
)

const certsEndpoint = "/now/certs"

// CertsClient contains the methods for the Cert API
type CertsClient struct {
	client *Client
}

// New creates a new cert
func (c CertsClient) New(domainNames []string) (Cert, ClientError) {
	params := CertParams{
		DomainNames: domainNames,
	}
	return c.NewFromParams(params)
}

// NewFromParams creates a new cert from params
func (c CertsClient) NewFromParams(params CertParams) (Cert, ClientError) {
	crt := Cert{}
	err := c.client.NewRequest("POST", certsEndpoint, params, &crt, nil)
	return crt, err
}

// Renew renews and existing cert
func (c CertsClient) Renew(domainNames []string) (Cert, ClientError) {
	crt := Cert{}
	params := CertParams{
		DomainNames: domainNames,
		Renew:       true,
	}
	err := c.client.NewRequest("POST", certsEndpoint, params, &crt, nil)
	return crt, err
}

// CertParams contains all fields for create
type CertParams struct {
	DomainNames []string `json:"domains"`
	Renew       bool     `json:"renew"`
}

// List retrieves a list of all the domains under the account
func (c CertsClient) List() ([]*Cert, ClientError) {
	crt := &certListResponse{}
	err := c.client.NewRequest("GET", certsEndpoint, nil, crt, nil)
	return crt.Certs, err
}

type certListResponse struct {
	Certs []*Cert `json:"certificates"`
}

// Delete deletes the domain by its ID
func (c CertsClient) Delete(domainName string) ClientError {
	return c.client.NewRequest("DELETE", fmt.Sprintf("%s/%s", certsEndpoint, domainName), nil, nil, nil)
}
