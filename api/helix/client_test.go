package helix

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mock_helix "github.com/c-rainbow/simplechatbot/api/helix/mock"
	"github.com/golang/mock/gomock"
	helix_api "github.com/nicklaw5/helix"
)

func TestEmpty(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mock_helix.NewMockInnerClientT(controller)
}

func TestGetUsers(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	testIds := []string{"111", "222", "333"}
	testUsernames := []string{"user111", "user222", "user333"}
	testResponse := &helix_api.UsersResponse{
		Data: helix_api.ManyUsers{
			Users: []helix_api.User{
				helix_api.User{ID: "111", Login: "user111"},
				helix_api.User{ID: "222", Login: "user222"},
			},
		},
	}

	innerClient := mock_helix.NewMockInnerClientT(controller)
	innerClient.EXPECT().GetUsers(&helix_api.UsersParams{
		IDs:    testIds,
		Logins: testUsernames,
	}).Return(testResponse, nil)

	client := &HelixClient{clientID: "12345", innerClient: innerClient}
	users, err := client.GetUsers([]string{"111", "222", "333"}, []string{"user111", "user222", "user333"})

	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))
	assert.Equal(t, users[0].ID, "111")
	assert.Equal(t, users[0].Login, "user111")
	assert.Equal(t, users[1].ID, "222")
	assert.Equal(t, users[1].Login, "user222")

}
