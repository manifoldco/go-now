package now

const planEndpoint = "/plan"

// PlansClient contains the methods for the Plan API
type PlansClient struct {
	client *Client
}

// Current returns the authenticated user's subscription
func (c PlansClient) Current() (Subscription, ClientError) {
	r := planResponse{}
	err := c.client.NewRequest("GET", planEndpoint, nil, &r, nil)
	return r.Subscription, err
}

type planResponse struct {
	Subscription Subscription `json:"subscription,omitempty"`
}

// TODO: implement Set
