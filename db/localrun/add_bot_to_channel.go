package localrun

import (
	"fmt"

	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/repository"
)

// AddBotToChannel Add new bot to channel
func AddBotToChannel() {
	repo := repository.NewBaseRepository()

	existingChannels := repo.GetAllChannels()
	fmt.Println("All existing channels: ", existingChannels[0])

	// in CreateChannels function, the created channel has BotIDs {111, 222, 333}
	botInfo := &models.Bot{TwitchID: DefaultBotTwitchID}
	// The channel already has bot 333
	chanInfo := &models.Channel{
		TwitchID: DefaultChannelTwitchID,
		Username: DefaultChannelUsername,
	}
	err := repo.AddBotToChannel(botInfo, chanInfo)
	if err != nil {
		fmt.Println("Error while adding bot ", botInfo.TwitchID, "to channel: ", err.Error())
	}
	channels := repo.GetAllChannels()
	fmt.Println("Channel status: ", channels[0]) // Should have bot 333, and not 444

	// Add botID that is not associated with the channel yet
	/*newBotInfo := &models.Bot{TwitchID: 444}
	err = repo.AddBotToChannel(newBotInfo, chanInfo)
	if err != nil {
		fmt.Println("Error while adding bot ", newBotInfo.TwitchID, "to channel: ", err.Error())
	}
	channels = repo.GetAllChannels()
	fmt.Println("Channel status: ", channels[0]) // Should have bot 333, and not 444
	*/
}
