package localrun

import (
	"fmt"

	"github.com/c-rainbow/simplechatbot/repository"
)

// Add new bot to local DynamoDB
func GetChannelsForBot() {
	repo := repository.NewBaseRepository()

	// in CreateChannels function, the created channel has BotIDs {111, 222, 333}
	botID := int64(222)
	channels := repo.GetAllChannelsForBot(botID)
	fmt.Println("Number of channels for bot ", botID, ": ", len(channels))

	// Test with bot with no associated channel
	botID = int64(555)
	channels = repo.GetAllChannelsForBot(botID)
	fmt.Println("Number of channels for bot ", botID, ": ", len(channels))
}
