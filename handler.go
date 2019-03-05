package simplechatbot

import (
	"strings"

	"github.com/lrstanley/girc"
)

type MessageHandler struct {
	// commands to handle
	userCommands UserCommands
}

func (handler *MessageHandler) handlePrivmsg(c *girc.Client, e girc.Event) {

	if e.EmptyTrailing {
		return
	}

	splited := strings.SplitN(e.Trailing, " ", 2)
	commandName := splited[0]
	username := e.Source.Name //

	if commandName == "!안녕" {
		c.Cmd.Reply(e, "안녕하세요")
	}

	command := handler.getCommand(username, commandName)

	if command.Name != "" {
		c.Cmd.Reply(e, "user: "+username) // command.response)
	}
}

func (handler *MessageHandler) getCommand(username, commandName string) Command {
	for _, userCommand := range handler.userCommands {
		if userCommand.User.Username == username {
			for _, command := range userCommand.Commands {
				if command.Name == commandName {
					return command
				}
			}
		}
	}
	return Command{}
}
