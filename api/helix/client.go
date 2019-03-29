package helix

import (
	"log"
	"strconv"

	helix_api "github.com/nicklaw5/helix"
)

const (
	DefaultClientID = "lc4tcxdkp0hkg87merghpp1f52alaj"
)

type HelixClientT interface {
	GetUsers(ids []int64, usernames []string) ([]helix_api.User, error)
}

type HelixClient struct {
	clientID    string
	innerClient *helix_api.Client
}

var _ HelixClientT = (*HelixClient)(nil)

func DefaultHelixClient() HelixClientT {
	return NewHelixClient(DefaultClientID)
}

func NewHelixClient(clientID string) HelixClientT {
	client, err := helix_api.NewClient(&helix_api.Options{
		ClientID: clientID,
	})
	if err != nil {
		// handle error
		log.Println("Error with Helix client", err.Error())
		return nil
	}

	return &HelixClient{clientID: clientID, innerClient: client}
}

func (client *HelixClient) GetUsers(ids []int64, usernames []string) ([]helix_api.User, error) {
	// Convert int64 ids to string
	var stringIDs []string
	if ids != nil {
		stringIDs := make([]string, len(ids))
		for index, id64 := range ids {
			stringIDs[index] = strconv.FormatInt(id64, 10)
		}
	}

	resp, err := client.innerClient.GetUsers(&helix_api.UsersParams{
		IDs:    stringIDs,
		Logins: usernames,
	})
	if err != nil {
		log.Println("Error while getting users", err.Error())
		return nil, err
	}
	return resp.Data.Users, nil
}

/*
func Client() {

	resp, err := client.GetUsers(&helix_api.UsersParams{
		IDs:    []string{"18074328"},
		Logins: []string{"c_rainbow"},
	})
	if err != nil {
		// handle error
	}

	fmt.Printf("Status code: %d\n", resp.StatusCode)
	fmt.Printf("Rate limit: %d\n", resp.GetRateLimit())
	fmt.Printf("Rate limit remaining: %d\n", resp.GetRateLimitRemaining())
	fmt.Printf("Rate limit reset: %d\n\n", resp.GetRateLimitReset())

	for _, user := range resp.Data.Users {
		fmt.Printf("ID: %s Name: %s\n", user.ID, user.DisplayName)
		fmt.Printf("Description: %s Email: %s\n", user.Description, user.Email)
		// fmt.Printf("ID: %s Name: %s\n", user.ID, user.DisplayName)
	}
}
*/
