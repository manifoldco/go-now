package now

// PlansClient contains the methods for the Plan API
type PlansClient struct {
	client *Client
}

// Current returns the authenticated user's subscription
func (c PlansClient) Current() (Subscription, error) {
	r := planResponse{}
	err := c.client.NewRequest("GET", "/plan", nil, &r)
	return r.Subscription, err
}

type planResponse struct {
	Subscription Subscription `json:"subscription,omitempty"`
}

// TODO: implement Set
