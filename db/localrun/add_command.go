package localrun

import (
	"fmt"

	"github.com/c-rainbow/simplechatbot/commands"

	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
	"github.com/c-rainbow/simplechatbot/plugins/chat/commandplugins"
	"github.com/c-rainbow/simplechatbot/repository"
)

// Add new bot to channel
func AddNewCommand() {
	repo := repository.NewBaseRepository()

	// existingChannels := repo.GetAllChannels()
	// fmt.Println("All existing channels: ", existingChannels[0])

	// in CreateChannels function, the created channel has BotIDs {111, 222, 333}
	botInfo := &models.Bot{TwitchID: DefaultBotTwitchID}
	// The channel already has bot 333
	chanInfo := &models.Channel{
		TwitchID: DefaultChannelTwitchID,
		Username: DefaultChannelUsername,
	}

	responseMap := make(map[string]parser.ParsedResponse)
	responseMap[commands.DefaultResponseKey] = *parser.ParseResponse("@$(user) 명령어를 성공적으로 추가하였습니다")

	commandToAdd := models.Command{
		BotID:          botInfo.TwitchID,
		ChannelID:      chanInfo.TwitchID,
		Name:           "!addcom",
		PluginType:     commandplugins.AddCommandPluginType,
		Responses:      responseMap, // make(map[string]parser.ParsedResponse),
		CooldownSecond: 5,
		Enabled:        true,
		Group:          "",
	}
	err := repo.AddCommand(chanInfo.Username, &commandToAdd)
	if err != nil {
		fmt.Println("Error while adding bot ", botInfo.TwitchID, "to channel: ", err.Error())
	}
}
