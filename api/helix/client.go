package helix

import (
	"log"

	helix_api "github.com/nicklaw5/helix"
)

const (
	DefaultClientID = "lc4tcxdkp0hkg87merghpp1f52alaj"
)

type HelixClientT interface {
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

type HelixClient struct {
	clientID    string
	innerClient InnerClientT
}

var _ HelixClientT = (*HelixClient)(nil)
var _ InnerClientT = (*helix_api.Client)(nil)

func DefaultHelixClient() HelixClientT {
	return NewHelixClient(DefaultClientID)
}

func NewHelixClient(clientID string) HelixClientT {
	client, err := helix_api.NewClient(&helix_api.Options{
		ClientID: clientID,
	})
	if err != nil {
		log.Println("Error creating new Helix client", err.Error())
		return nil
	}

	return &HelixClient{clientID: clientID, innerClient: client}
}

// ids are numeric IDs.
func (client *HelixClient) GetUsers(ids []string, usernames []string) ([]helix_api.User, error) {
	// Convert int64 ids to string
	resp, err := client.innerClient.GetUsers(&helix_api.UsersParams{
		IDs:    ids,
		Logins: usernames,
	})
	if err != nil {
		log.Println("Error while getting users", err.Error())
		return nil, err
	}
	return resp.Data.Users, nil
}

// Here, fromID and toID are numeric IDs.
func (client *HelixClient) GetUsersFollows(fromID string, toID string) ([]helix_api.UserFollow, error) {
	resp, err := client.innerClient.GetUsersFollows(&helix_api.UsersFollowsParams{
		FromID: fromID, ToID: toID,
	})

	if err != nil {
		log.Println("Error while getting users follows", err.Error())
		return nil, err
	}
	return resp.Data.Follows, nil
}

// ids are numeric IDs.
func (client *HelixClient) GetStreams(ids []string, usernames []string) ([]helix_api.Stream, error) {
	resp, err := client.innerClient.GetStreams(&helix_api.StreamsParams{
		UserIDs: ids, UserLogins: usernames,
	})

	if err != nil {
		log.Println("Error while getting streams", err.Error())
		return nil, err
	}
	return resp.Data.Streams, nil
}
