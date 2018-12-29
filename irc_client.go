package simplechatbot

import (
	"os"

	"github.com/lrstanley/girc"
)

// Client is Wrapper struct for IRC client.
type Client struct {
	// Inner girc client
	*girc.Client
}

// NewClient creates new client.
func NewClient(username, oauthToken string) *Client {
	c := &Client{
		girc.New(girc.Config{
			Server:     TwitchIrcHost,
			Port:       TwitchIrcPort,
			Nick:       username, // Both Nick and User must be specified.
			User:       username,
			ServerPass: oauthToken,
			Debug:      os.Stdout,
		}),
	}

	return c
}
