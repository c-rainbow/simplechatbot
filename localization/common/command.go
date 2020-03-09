package common

// BotCommandLocaleConfig locale config for fallback bot responses
type BotCommandLocaleConfig struct {
	AddCommand    AddCommandMessages
	EditCommand   EditCommandMessages
	DeleteCommand DeleteCommandMessages
	ListCommands  ListCommandsMessages
	Uptime        UptimeMessages
}

// AddCommandMessages default messages for AddCommandPlugin when adding bots to channels
type AddCommandMessages struct {
	DefaultSuccess             string
	DefaultFailure             string
	NoPermission               string
	NotEnoughArgument          string
	TargetCommandAlreadyExists string
}

// EditCommandMessages default messages for EditCommandPlugin when adding bots to channels
type EditCommandMessages struct {
	DefaultSuccess         string
	DefaultFailure         string
	NoPermission           string
	NotEnoughArgument      string
	TargetCommandNotExists string
}

// DeleteCommandMessages default messages for DeleteCommandPlugin when adding bots to channels
type DeleteCommandMessages struct {
	DefaultSuccess         string
	DefaultFailure         string
	NoPermission           string
	NotEnoughArgument      string
	TargetCommandNotExists string
}

// ListCommandsMessages default messages for ListCommandsPlugin when adding bots to channels
type ListCommandsMessages struct {
	DefaultSuccess string
	DefaultFailure string
}

// UptimeMessages default messages for UptimePlugin when adding bots to channels
type UptimeMessages struct {
	DefaultOnline  string
	DefaultOffline string
}
