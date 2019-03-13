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
		TwitchID: DefaultChannelTwitchID,
		Username: DefaultChannelUsername,
	}
	repo.CreateNewChannel(testChannel)
	channels = repo.GetAllChannels()
	fmt.Println("New number of existing channels: ", len(channels))
	fmt.Println("Channel info: ", channels[0])
}
