package commands

import (
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
)

func BuildCommand(
	botID int64, channelID int64, name string, pluginType string,
	defaultResponse string) (*models.Command, error) {

	responseMap := make(map[string]models.ParsedResponse)
	parsedDefaultResponse := parser.ParseResponse(defaultResponse)
	err := parser.Validate(parsedDefaultResponse)
	if err != nil {
		return nil, err
	}

	responseMap[models.DefaultResponseKey] = *parsedDefaultResponse

	return &models.Command{
		BotID:          botID,
		ChannelID:      channelID,
		Name:           name,
		PluginType:     pluginType,
		Responses:      responseMap,
		CooldownSecond: 5,
		Permission:     models.PermissionEveryone,
		Enabled:        true,
		Group:          "",
	}, nil
}
