// Twitch IRC client module

package simplechatbot

import (
	twitch_irc "github.com/gempir/go-twitch-irc"
)

// TwitchClient is Wrapper struct for existing Twitch IRC client.
type TwitchClient struct {
	// Inner twitch IRC chat client
	*twitch_irc.Client
}

// NewTwitchClient creates new client.
func NewTwitchClient(username, oauthToken string) *TwitchClient {
	return &TwitchClient{
		twitch_irc.NewClient(username, oauthToken),
	}
}
