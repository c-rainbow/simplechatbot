package localrun

import (
	"fmt"

	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
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

	commandToAdd := BuildCommand(
		botInfo.TwitchID,
		chanInfo.TwitchID,
		"!addcom",
		commandplugins.AddCommandPluginType,
		"@$(user) 명령어를 성공적으로 추가하였습니다")
	err := repo.AddCommand(chanInfo.Username, commandToAdd)
	if err != nil {
		fmt.Println("Error while adding bot ", botInfo.TwitchID, "to channel: ", err.Error())
	}
}

func BuildCommand(
	botID int64, channelID int64, name string, pluginType string, defaultResponse string) *models.Command {

	responseMap := make(map[string]parser.ParsedResponse)
	responseMap[chatplugins.DefaultResponseKey] = *parser.ParseResponse(defaultResponse)

	return &models.Command{
		BotID:          botID,
		ChannelID:      channelID,
		Name:           name,
		PluginType:     pluginType,
		Responses:      responseMap,
		CooldownSecond: 5,
		Enabled:        true,
		Group:          "",
	}
}
