package mock

// Client is a mock implementation of linodego.Client
type Client struct{}

// NewClient is a mock implementation of linodego.Client constructor
func NewClient() *Client {
	return &Client{}
}
