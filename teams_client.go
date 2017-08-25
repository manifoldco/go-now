package now

import (
	"fmt"
)

const teamsEndpoint = "/teams"

// TeamsClient contains the methods for the Team API
type TeamsClient struct {
	client *Client
}

// New creates a new Team
func (c TeamsClient) New(slug string) (Team, ClientError) {
	return c.NewFromParams(TeamParams{
		Slug: slug,
	})
}

// NewFromParams creates a new Team from params
func (c TeamsClient) NewFromParams(params TeamParams) (Team, ClientError) {
	t := Team{}
	err := c.client.NewRequest("POST", teamsEndpoint, params, &t, nil)
	return t, err
}

// TeamParams contains possible fields for creation
type TeamParams struct {
	Slug string `json:"slug"`
}

// List retrieves a list of all the domains under the account
func (c TeamsClient) List() ([]Team, ClientError) {
	d := &teamListResponse{}
	err := c.client.NewRequest("GET", teamsEndpoint, nil, d, nil)
	return d.Teams, err
}

type teamListResponse struct {
	Teams []Team `json:"teams"`
}

// Members retrieves all members associated with a team
func (c TeamsClient) Members(teamID string) ([]TeamMember, ClientError) {
	var d []TeamMember
	err := c.client.NewRequest("GET", fmt.Sprintf("%s/%s/members", teamsEndpoint, teamID), nil, &d, nil)
	return d, err
}

// Delete deletes the domain by its ID
func (c TeamsClient) Delete(teamID string) ClientError {
	return c.client.NewRequest("DELETE", fmt.Sprintf("%s/%s", teamsEndpoint, teamID), nil, nil, nil)
}

// Rename updates the name value for the specified team
func (c TeamsClient) Rename(teamID, name string) ClientError {
	return c.client.NewRequest("POST", fmt.Sprintf("%s/%s/members", teamsEndpoint, teamID), &RenameParams{
		Name: name,
	}, nil, nil)
}

// RenameParams contains possible fields for rename
type RenameParams struct {
	Name string `json:"name"`
}

// InviteUser sends an invite for the specified team to the email provided
func (c TeamsClient) InviteUser(teamID, email string) ClientError {
	return c.client.NewRequest("POST", fmt.Sprintf("%s/%s/members", teamsEndpoint, teamID), &InviteParams{
		Email: email,
	}, nil, nil)
}

// InviteParams contains possible fields for invite
type InviteParams struct {
	Email string `json:"email"`
}
