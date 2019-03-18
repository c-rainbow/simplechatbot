package simplechatbot

import (
	"fmt"
	"log"
	"sort"
	"strings"

	commands "github.com/c-rainbow/simplechatbot/commands"
	models "github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

var (
	ReservedCommands = []string{"!addcom", "!editcom", "!delcom"}
)

type ChatMessageHandlerT interface {
	OnNewMessage(channel string, sender twitch_irc.User, message twitch_irc.Message)
}

type ChatMessageHandler struct {
	botInfo   *models.Bot
	repo      SingleBotRepositoryT
	ircClient *TwitchClient
	chatters  map[string]bool
}

var _ ChatMessageHandlerT = (*ChatMessageHandler)(nil)

func NewChatMessageHandler(
	botInfo *models.Bot, repo SingleBotRepositoryT, ircClient *TwitchClient) *ChatMessageHandler {
	return &ChatMessageHandler{botInfo: botInfo, repo: repo, ircClient: ircClient, chatters: make(map[string]bool)}
}

func (handler *ChatMessageHandler) OnNewMessage(
	channel string, sender twitch_irc.User, message twitch_irc.Message) {
	fmt.Println("Chat received: ", message.Raw)
	commandName := getCommandName(message.Text)
	commandName = strings.ToLower(commandName)
	toSay := ""

	// TODO: Which bot should react to this?
	if commands.ReservedCommands[commandName] {
		var err error
		if commandName == commands.ListCommandsKey || commandName == commands.ListCommandsKeyKor {
			var commands []*models.Command
			commands, err = handler.repo.ListCommands(channel)
			commandNames := make([]string, len(commands))
			for i, command := range commands {
				fmt.Println("Command name: ", command.Name)
				commandNames[i] = command.Name
			}
			fmt.Println("command names:", commandNames)
			sort.Strings(commandNames)
			toSay = strings.Join(commandNames, ", ")
			fmt.Println("toSay:", toSay)
		} else {

			parsedCommand, err := commands.ParseCommand(handler.botInfo.TwitchID, message.Text, channel, &sender, &message)
			if err != nil {
				log.Println("ParseCommand err: ", err.Error())
				return
			}
			switch commandName {
			case commands.AddCommandKey, commands.AddCommandKeyKor:
				err = handler.repo.AddCommand(channel, parsedCommand)
			case commands.EditCommandKey, commands.EditCommandKeyKor:
				err = handler.repo.EditCommand(channel, parsedCommand)
			case commands.DeleteCommandKey, commands.DeleteCommandKeyKor:
				err = handler.repo.DeleteCommand(channel, parsedCommand)

			default:
				log.Println("Cannot find commandName", commandName)
				return
			}
		}
		if err != nil {
			log.Println("Failed to process command update: ", err.Error())
			return
		}
	} else {
		// Check if command with the same name exists
		command := handler.repo.GetCommandByChannelAndName(channel, commandName)
		// displayName := sender.DisplayName

		if command != nil {
			// TODO: diverse response cases
			response := command.Responses[commands.DefaultResponseKey]
			converted, err := parser.ConvertResponse(&response, channel, &sender, &message)
			if err != nil {
				fmt.Println("ERror while converting response: ", err.Error())
				return
			}
			toSay = converted
		}

		/*if commandName == "hello" || commandName == "hi" {
			toSay = "Hi " + displayName
		} else if commandName == "안녕하세요" && sender.Username != "c_rainbow" {
			toSay = displayName + " 님도 안녕하세요"
		} else */
		if commandName == "!quit" && sender.Username == "c_rainbow" {
			handler.ircClient.Depart(channel)
			handler.ircClient.Disconnect()
		}

		// Check for new chatter
		/*if _, has := handler.chatters[displayName]; !has {
			handler.chatters[displayName] = true
			toSay = displayName + " 님 어서오세요 환영합니다"
		}*/
	}

	if toSay != "" {
		handler.ircClient.Say(channel, toSay)
	}
}

// Gets command name from the full chat text
func getCommandName(text string) string {
	// TODO: Is there always to heading space?
	index := strings.Index(text, " ")
	// If there is no space in the chat text, then the chat itself is the command
	if index == -1 {
		return text
	}
	return text[:index]
}
