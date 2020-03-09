package americanenglish

import "github.com/c-rainbow/simplechatbot/localization/common"

var (
	addCommandMessages = common.AddCommandMessages{
		DefaultSuccess:             "@$(user) 새 명령어를 추가하였습니다",
		DefaultFailure:             "실패하였습니다",
		NoPermission:               "@$(user) 추가 권한이 없습니다",
		NotEnoughArgument:          "입력 형식이 잘못되었습니다",
		TargetCommandAlreadyExists: "같은 명령어가 이미 있습니다",
	}

	editCommandMessages = common.EditCommandMessages{
		DefaultSuccess:         "@$(user) 명령어를 수정하였습니다",
		DefaultFailure:         "명령어를 수정하지 못했습니다",
		NoPermission:           "@$(user) 수정 권한이 없습니다",
		NotEnoughArgument:      "입력 형식이 잘못되었습니다",
		TargetCommandNotExists: "수정할 명령어가 존재하지 않습니다",
	}

	deleteCommandMessages = common.DeleteCommandMessages{
		DefaultSuccess:         "@$(user) 명령어를 삭제하였습니다",
		DefaultFailure:         "명령어를 삭제하지 못했습니다",
		NoPermission:           "@$(user) 삭제 권한이 없습니다",
		NotEnoughArgument:      "입력 형식이 잘못되었습니다",
		TargetCommandNotExists: "삭제할 명령어가 존재하지 않습니다",
	}

	listCommandsMessages = common.ListCommandsMessages{
		DefaultSuccess: "$(args0)",
		DefaultFailure: "명령어를 가져올 수 없습니다",
	}

	uptimeMessages = common.UptimeMessages{
		DefaultOnline:  "$(uptime)동안 방송중입니다",
		DefaultOffline: "지금은 방송중이 아닙니다",
	}

	botCommandLocale = common.BotCommandLocaleConfig{
		AddCommand:    addCommandMessages,
		EditCommand:   editCommandMessages,
		DeleteCommand: deleteCommandMessages,
		ListCommands:  listCommandsMessages,
		Uptime:        uptimeMessages,
	}
)
