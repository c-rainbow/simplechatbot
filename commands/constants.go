package commands

// Reserved commands. They cannot be added/edited/deleted
var (
	AddCommandKey       = "!addcom"
	EditCommandKey      = "!editcom"
	DeleteCommandKey    = "!delcom"
	ListCommandsKey     = "!commands"
	AddCommandKeyKor    = "!추가"
	EditCommandKeyKor   = "!수정"
	DeleteCommandKeyKor = "!삭제"
	ListCommandsKeyKor  = "!명령어"
)

var (
	// Default response key in resopnse map
	DefaultResponseKey = "DEFAULT"
)

// All reserved commands should be added here
var ReservedCommands = map[string]bool{
	AddCommandKey:       true,
	EditCommandKey:      true,
	DeleteCommandKey:    true,
	ListCommandsKey:     true,
	AddCommandKeyKor:    true,
	EditCommandKeyKor:   true,
	DeleteCommandKeyKor: true,
	ListCommandsKeyKor:  true,
}
