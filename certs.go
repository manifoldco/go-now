package now

import "time"

// Cert is the contents of an ssl certificate object
type Cert struct {
	UID     string     `json:"uid"`
	Created *time.Time `json:"created,omitempty"`
}
