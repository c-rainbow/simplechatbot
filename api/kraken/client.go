package helix

import (
	helix_api "github.com/nicklaw5/helix"
)

const (
	DefaultClientID = "lc4tcxdkp0hkg87merghpp1f52alaj"
)

type KrakenClientT interface {
	GetUsers(ids []string, usernames []string) ([]helix_api.User, error)
	GetUsersFollows(fromID string, toID string) ([]helix_api.UserFollow, error)
	GetStreams(ids []string, usernames []string) ([]helix_api.Stream, error)
}

// InnerClientT Interface for inner API client.
type InnerClientT interface {
	GetUsers(params *helix_api.UsersParams) (*helix_api.UsersResponse, error)
	GetUsersFollows(params *helix_api.UsersFollowsParams) (*helix_api.UsersFollowsResponse, error)
	GetStreams(params *helix_api.StreamsParams) (*helix_api.StreamsResponse, error)
}
