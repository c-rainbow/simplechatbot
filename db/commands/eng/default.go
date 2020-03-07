package eng

import (
	"log"

	"github.com/c-rainbow/simplechatbot/db/commands"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/plugins/chat/commandplugins"
	"github.com/c-rainbow/simplechatbot/plugins/chat/selfban"
)

func DefaultAddCommandModel(channel *models.Channel, botID int64) {
	commandName := "!addcom"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, commandplugins.AddCommandPluginType,
		"@$(user) Successfully added a new command")
	if err != nil {
		log.Println("Failed to create a new command:", commandName, err.Error())
		return
	}
	channel.Commands[commandName] = *command
}

func DefaultEditCommandModel(channel *models.Channel, botID int64) {
	commandName := "!editcom"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, commandplugins.EditCommandPluginType,
		"@$(user) Successfully edited the command")
	if err != nil {
		log.Println("Failed to create a new command:", commandName, err.Error())
		return
	}
	channel.Commands[commandName] = *command
}

func DefaultDeleteCommandModel(channel *models.Channel, botID int64) {
	commandName := "!delcom"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, commandplugins.DeleteCommandPluginType,
		"@$(user) Successfully deleted the command")
	if err != nil {
		log.Println("Failed to create a new command:", commandName, err.Error())
		return
	}
	channel.Commands[commandName] = *command
}

func DefaultListCommandsModel(channel *models.Channel, botID int64) {
	commandName := "!commands"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, commandplugins.ListCommandsPluginType,
		"The list of commands are: $(arg0)")
	if err != nil {
		log.Println("Failed to create a new command:", commandName, err.Error())
		return
	}
	channel.Commands[commandName] = *command
}

func DefaultSelfBanGameModel(channel *models.Channel, botID int64) {
	commandName := "!banme"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, selfban.BanOneselfPluginType, "@$(user) You are banned for 3 seconds")
	if err != nil {
		log.Println("Failed to create a new command:", commandName, err.Error())
		return
	}
	channel.Commands[commandName] = *command
}
