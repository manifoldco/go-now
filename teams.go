package now

import "time"

// Team is the contents of a team object
type Team struct {
	ID        string     `json:"id"`
	Slug      string     `json:"slug"`
	Name      string     `json:"name"`
	CreatorID string     `json:"creator_id"`
	Created   *time.Time `json:"created,omitempty"`
}

// TeamMember represents a membership to a team
type TeamMember struct {
	UID      string `json:"uid"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
