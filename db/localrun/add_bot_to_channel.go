package localrun

import (
	"fmt"

	simplechatbot "github.com/c-rainbow/simplechatbot"
	"github.com/c-rainbow/simplechatbot/models"
)

// Add new bot to channel
func AddBotToChannel() {
	repo := simplechatbot.NewBaseRepository()

	existingChannels := repo.GetAllChannels()
	fmt.Println("All existing channels: ", existingChannels[0])

	// in CreateChannels function, the created channel has BotIDs {111, 222, 333}
	botInfo := &models.Bot{TwitchID: 333}
	// The channel already has bot 333
	chanInfo := &models.Channel{
		TwitchID:    1234,
		Username:    "test_channel",
		DisplayName: "Testchannel",
		BotIDs:      []int64{111, 222, 333},
		Commands: []models.Command{
			models.Command{
				UUID:      "a1b2c3d4e5f6",
				BotID:     222,
				Name:      "!test",
				ChannelID: 1234,
			},
		},
	}
	err := repo.AddBotToChannel(botInfo, chanInfo)
	if err != nil {
		fmt.Println("Error while adding bot ", botInfo.TwitchID, "to channel: ", err.Error())
	}
	channels := repo.GetAllChannels()
	fmt.Println("Channel status: ", channels[0]) // Should have bot 333, and not 444

	// Add botID that is not associated with the channel yet
	newBotInfo := &models.Bot{TwitchID: 444}
	err = repo.AddBotToChannel(newBotInfo, chanInfo)
	if err != nil {
		fmt.Println("Error while adding bot ", newBotInfo.TwitchID, "to channel: ", err.Error())
	}
	channels = repo.GetAllChannels()
	fmt.Println("Channel status: ", channels[0]) // Should have bot 333, and not 444

}
