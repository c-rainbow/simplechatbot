package kor

import (
	"log"

	"github.com/c-rainbow/simplechatbot/db/commands"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/plugins/chat/commandplugins"
	"github.com/c-rainbow/simplechatbot/plugins/chat/games"
	"github.com/c-rainbow/simplechatbot/plugins/chat/selfban"
)

func DefaultAddCommandModel(channel *models.Channel, botID int64) {
	commandName := "!추가"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, commandplugins.AddCommandPluginType,
		"@$(user) 명령어를 성공적으로 추가하였습니다")
	if err != nil {
		log.Println("명령어 '!추가'를 추가하는데 실패하였습니다", err.Error())
		return
	}
	channel.Commands[commandName] = *command
}

func DefaultEditCommandModel(channel *models.Channel, botID int64) {
	commandName := "!수정"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, commandplugins.EditCommandPluginType,
		"@$(user) 명령어를 성공적으로 수정하였습니다")
	if err != nil {
		log.Println("명령어 '!수정'을 추가하는데 실패하였습니다", err.Error())
		return
	}
	channel.Commands[commandName] = *command
}

func DefaultDeleteCommandModel(channel *models.Channel, botID int64) {
	commandName := "!삭제"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, commandplugins.DeleteCommandPluginType,
		"@$(user) 명령어를 성공적으로 삭제하였습니다")
	if err != nil {
		log.Println("명령어 '!삭제'를 추가하는데 실패하였습니다", err.Error())
		return
	}
	channel.Commands[commandName] = *command
}

func DefaultListCommandsModel(channel *models.Channel, botID int64) {
	commandName := "!명령어"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, commandplugins.ListCommandsPluginType,
		"명령어 모음: $(arg0)")
	if err != nil {
		log.Println("명령어 '!명령어'를 추가하는데 실패하였습니다", err.Error())
		return
	}
	channel.Commands[commandName] = *command
}

func DefaultNumberGuesserModel(channel *models.Channel, botID int64) {
	commandName := "!숫자"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, games.NumberGuesserPluginType, ".")
	if err != nil {
		log.Println("명령어 '!숫자'를 추가하는데 실패하였습니다", err.Error())
		return
	}
	channel.Commands[commandName] = *command
}

func DefaultDiceGameModel(channel *models.Channel, botID int64) {
	commandName := "!주사위"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, games.DicePluginType, ".")
	if err != nil {
		log.Println("명령어 '!주사위'를 추가하는데 실패하였습니다", err.Error())
		return
	}
	channel.Commands[commandName] = *command
}

func DefaultSelfBanGameModel(channel *models.Channel, botID int64) {
	commandName := "!셀프밴"
	command, err := commands.BuildCommand(
		botID, channel.TwitchID, commandName, selfban.BanOneselfPlugin, "@$(user) 님 밴")
	if err != nil {
		log.Println("명령어 '!셀프밴'을 추가하는데 실패하였습니다", err.Error())
		return
	}
	channel.Commands[commandName] = *command
}
