package now

import "time"

// Domain is the contents of a domain object
type Domain struct {
	UID         string     `json:"uid"`
	Verified    bool       `json:"verified"`
	VerifyToken string     `json:"verifyToken,omitempty"`
	Created     *time.Time `json:"created,omitempty"`
}
