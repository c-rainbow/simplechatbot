package americanenglish

import "github.com/c-rainbow/simplechatbot/localization/common"

var (
	addCommandMessages = common.AddCommandMessages{
		DefaultSuccess:             "@$(user) successfully added the new command",
		DefaultFailure:             "Failed to add the new command",
		NoPermission:               "@$(user) You don't have permission to add command",
		NotEnoughArgument:          "Invalid format",
		TargetCommandAlreadyExists: "The command already exists",
	}

	editCommandMessages = common.EditCommandMessages{
		DefaultSuccess:         "@$(user) Successfully edited the command",
		DefaultFailure:         "Failed to edit the command",
		NoPermission:           "@$(user) You don't have permission to edit command",
		NotEnoughArgument:      "Invalid format",
		TargetCommandNotExists: "Command does not exist",
	}

	deleteCommandMessages = common.DeleteCommandMessages{
		DefaultSuccess:         "@$(user) The command is now deleted",
		DefaultFailure:         "Could not delete the command",
		NoPermission:           "@$(user) You don't have permission to delete command",
		NotEnoughArgument:      "Invalid format",
		TargetCommandNotExists: "Command does not exist",
	}

	listCommandsMessages = common.ListCommandsMessages{
		DefaultSuccess: "$(args0)",
		DefaultFailure: "Cannot get the list of commands right now",
	}

	uptimeMessages = common.UptimeMessages{
		DefaultOnline:  "Uptime: $(uptime)",
		DefaultOffline: "The stream is offline",
	}

	botCommandLocale = common.BotCommandLocaleConfig{
		AddCommand:    addCommandMessages,
		EditCommand:   editCommandMessages,
		DeleteCommand: deleteCommandMessages,
		ListCommands:  listCommandsMessages,
		Uptime:        uptimeMessages,
	}
)
