// Data manager with LevelDB

package simplechatbot

import (
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
)

// PrototypeDataManager Data manager for prototype chatbot system.
// Single bot, single user. In-memory command
type PrototypeDataManager struct {
	commandMap     map[string]Command
	databaseHandle *leveldb.DB
}

func (manager *PrototypeDataManager) getDatabase() *leveldb.DB {
	if manager.databaseHandle != nil {
		return manager.databaseHandle
	}

	dbFile, _ := filepath.Abs("../db/database.db")
	manager.databaseHandle, _ = leveldb.OpenFile(dbFile, nil)
	return manager.databaseHandle
}

func (manager *PrototypeDataManager) getUserInfo() User {
	return User{Username: "c_rainbow"}
}

func (manager *PrototypeDataManager) getCommandsPerUser(user User) Commands {
	var commands []Command
	for _, v := range manager.commandMap {
		commands = append(commands, v)
	}
	return commands
}

func (manager *PrototypeDataManager) editCommand(user User, command Command) {
	manager.commandMap[command.Name] = command
}

func (manager *PrototypeDataManager) addCommand(user User, command Command) {
	_ = user
	manager.commandMap[command.Name] = command
}

func (manager *PrototypeDataManager) deleteCommand(user User, command Command) {
	delete(manager.commandMap, command.Name)
}

func (manager *PrototypeDataManager) getBotInfo() Bot {
	db := manager.getDatabase()

	botUsername, _ := db.Get([]byte(BotUsernameKey), nil)
	botOauthToken, _ := db.Get([]byte(BotOauthTokenKey), nil)

	// botUsername = "intern_bot"
	// botOauthToken

	_, _ = botUsername, botOauthToken

	db.Close()

	return Bot{
		Name: "intern_bot", OauthToken: "oauth:s1f67e8dn7gufsuq0jyssnucwj39sq",
	}
}
