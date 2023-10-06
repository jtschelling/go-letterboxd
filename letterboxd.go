package letterboxd

const (
	// APIBaseURL is used to create a client
	APIBaseURL = "https://api.letterboxd.com/api/v0"
)

// Client provides a common object and plain configuration for interacting with this library
type Client struct {
	BaseURL      string
	clientID     string
	clientSecret string
}

// New creates a letterboxd API client from an auth token
func New(clientID string, clientSecret string) *Client {
	return &Client{
		BaseURL:      APIBaseURL,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}
