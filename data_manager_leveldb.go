// Data manager with LevelDB

package simplechatbot

import (
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Bot auth related constants
const (
	BotUsernameKey   = "bot_info.username"
	BotOauthTokenKey = "bot_info.oauth_token"
	DatabaseFilePath = "./db/database.db" // DatabaseFilePath has DB file path.
)

// OwnerUsernameKey has bot owner Twitch username.
const OwnerUsernameKey = "owner_info.username"

const CommandKey = "command_info."

type LevelDbDataManager struct {
	databaseHandle *leveldb.DB
}

func (manager *LevelDbDataManager) getDatabase() *leveldb.DB {
	if manager.databaseHandle != nil {
		return manager.databaseHandle
	}

	dbFile, _ := filepath.Abs(DatabaseFilePath)
	manager.databaseHandle, _ = leveldb.OpenFile(dbFile, nil)
	return manager.databaseHandle
}

/*
	addCommand(user User, command Command)
	editCommand(user User, command Command)
	deleteCommand(user User, command Command)
*/
func (manager *LevelDbDataManager) getBotInfo() Bot {
	db := manager.getDatabase()

	botUsername, _ := db.Get([]byte(BotUsernameKey), nil)
	botOauthToken, _ := db.Get([]byte(BotOauthTokenKey), nil)

	return Bot{
		Name: string(botUsername), OauthToken: string(botOauthToken),
	}
}

func (manager *LevelDbDataManager) getUserInfo() User {
	db := manager.getDatabase()

	ownerUsername, _ := db.Get([]byte(OwnerUsernameKey), nil)
	return User{Username: string(ownerUsername)}
}

func (manager *LevelDbDataManager) getCommandsPerUser(user User) Commands {
	db := manager.getDatabase()
	commandKey := CommandKey + user.Username + "-"
	iter := db.NewIterator(util.BytesPrefix([]byte(commandKey)), nil)

	var commands Commands
	for iter.Next() {
		value := iter.Value()
		commands = append(commands, manager.UnmarshalCommand(value))
	}
	return commands
}

func (manager *LevelDbDataManager) addCommand(ownerUsername, name, response string, cooldownSecond int) {
	commandKey := CommandKey + ownerUsername + "." + ownerUsername
	command := Command{
		Name: name, Response: response, CooldownSecond: cooldownSecond,
	}

	db := manager.getDatabase()
	db.Put([]byte(commandKey), manager.MarshalCommand(&command), nil)

}

func (manager *LevelDbDataManager) UnmarshalCommand(data []byte) Command {
	c := Command{}
	json.Unmarshal(data, c)
	return c
}

func (manager *LevelDbDataManager) MarshalCommand(c *Command) []byte {
	bytes, err := json.Marshal(c)
	if err != nil {
		log.Fatal("Error marshalling JSON: " + err.Error())
	}
	return bytes
}
