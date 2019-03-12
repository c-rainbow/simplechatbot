package localrun

import (
	"fmt"

	simplechatbot "github.com/c-rainbow/simplechatbot"
	"github.com/c-rainbow/simplechatbot/models"
)

// Add new channel to local DynamoDB
func AddNewChannel() {
	repo := simplechatbot.NewBaseRepository()

	// Make sure that no channels exist first
	channels := repo.GetAllChannels()
	fmt.Println("Number of existing channel: ", len(channels))

	// Add one channel fixture
	testChannel := &models.Channel{
		TwitchID:    1234,
		Username:    "test_channel",
		DisplayName: "Testchannel",
		BotIDs:      []int64{111, 222, 333},
		Commands: []models.Command{
			models.Command{
				UUID:      "a1b2c3d4e5f6",
				BotID:     233,
				Name:      "!test",
				ChannelID: 1234,
			},
		},
	}
	repo.CreateNewChannel(testChannel)
	channels = repo.GetAllChannels()
	fmt.Println("New number of existing channels: ", len(channels))
	fmt.Println("Channel info: ", channels[0])
}
