package simplechatbot

// DataManager is responsible for converting persistent data to Go struct, and vice versa.
type DataManager interface {
	getBotInfo() Bot
	getUserInfo() User
	getCommandsPerUser(user User) Commands
	addCommand(user User, command Command)
	editCommand(user User, command Command)
	deleteCommand(user User, command Command)
}
