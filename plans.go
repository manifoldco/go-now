package now

// Subscription represents a user's subscription state
type Subscription struct {
	ID   string `json:"id"`
	Plan Plan   `json:"plan"`
	// TODO: Add remaining subscription fields
}

// Plan contains all fields relevant to a plan object
type Plan struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Interval      string `json:"interval"`
	IntervalCount int64  `json:"interval_count"`
}
