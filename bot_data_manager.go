package simplechatbot

// DataManager is responsible for converting persistent data to Go struct, and vice versa.
type BotDataManager interface {
	GetBotInfo() *Bot
	GetUserInfos() []*User
	GetCommandsPerUser(user *User) []*Command
	GetAllUserCommands() []*UserCommand
	AddCommand(user *User, command *Command)
	EditCommand(user *User, command *Command)
	DeleteCommand(user *User, command *Command)
}
