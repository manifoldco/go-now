package now

import (
	"fmt"
)

// TeamsClient contains the methods for the Team API
type TeamsClient struct {
	client *Client
}

// New creates a new Team
func (c TeamsClient) New(slug string) (Team, error) {
	d := Team{}
	err := c.client.NewRequest("POST", "/teams", teamParams{
		Slug: slug,
	}, &d)
	return d, err
}

type teamParams struct {
	Slug string `json:"slug"`
}

// List retrieves a list of all the domains under the account
func (c TeamsClient) List() ([]Team, error) {
	d := &teamListResponse{}
	err := c.client.NewRequest("GET", "/teams", nil, d)
	return d.Teams, err
}

type teamListResponse struct {
	Teams []Team `json:"teams"`
}

// Members retrieves all members associated with a team
func (c TeamsClient) Members(teamID string) ([]TeamMember, error) {
	var d []TeamMember
	err := c.client.NewRequest("GET", fmt.Sprintf("/teams/%s/members", teamID), nil, &d)
	return d, err
}

// Delete deletes the domain by its ID
func (c TeamsClient) Delete(teamID string) error {
	return c.client.NewRequest("DELETE", fmt.Sprintf("/teams/%s", teamID), nil, nil)
}

// Rename updates the name value for the specified team
func (c TeamsClient) Rename(teamID, name string) error {
	return c.client.NewRequest("POST", fmt.Sprintf("/teams/%s/members", teamID), &renameParams{
		Name: name,
	}, nil)
}

type renameParams struct {
	Name string `json:"name"`
}

// InviteUser sends an invite for the specified team to the email provided
func (c TeamsClient) InviteUser(teamID, email string) error {
	return c.client.NewRequest("POST", fmt.Sprintf("/teams/%s/members", teamID), &inviteParams{
		Email: email,
	}, nil)
}

type inviteParams struct {
	Email string `json:"email"`
}
