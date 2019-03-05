// Data manager with LevelDB

package simplechatbot

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Bot auth related constants
const (
	DefaultDbFilePath = "./db/data.db" // DatabaseFilePath has DB file path.
)

type SqLiteDataManager struct {
	dbFilePath     string
	databaseHandle *sql.DB
}

func NewSqLiteDataManager(dbFilePath string) *SqLiteDataManager {
	manager := &SqLiteDataManager{dbFilePath: dbFilePath}
	return manager
}

func (manager *SqLiteDataManager) getDatabase() *sql.DB {

	// Initialize DB handle only for the first time.
	if manager.databaseHandle == nil {
		db, err := sql.Open("sqlite3", manager.dbFilePath)
		if err != nil {
			log.Fatal(err)
		}

		manager.databaseHandle = db
		db.Ping() // Just to make sure connection works
	}
	return manager.databaseHandle
}

/*
	addCommand(user User, command Command)
	editCommand(user User, command Command)
	deleteCommand(user User, command Command)
*/
func (manager *SqLiteDataManager) getBotInfo() Bot {
	db := manager.getDatabase()

	var botUsername string
	var botOauthToken string

	rows, err := db.Query("SELECT username, oauth_token FROM Bots LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&botUsername, &botOauthToken)

	return Bot{Name: botUsername, OauthToken: botOauthToken}
}

func (manager *SqLiteDataManager) getUserInfo() User {
	db := manager.getDatabase()

	var username string

	rows, err := db.Query("SELECT username FROM Users LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	rows.Next()
	rows.Scan(&username)

	return User{Username: username}
}

func (manager *SqLiteDataManager) getCommandsPerUser(user User) Commands {
	sqlText := `
	SELECT c.name AS command_name, r.name AS response_name
	FROM Commands c
	INNER JOIN Responses r
	ON r.id = c.response_id
	INNER JOIN Users u
	ON u.id = c.user_id
	WHERE u.username = ` + user.Username

	db := manager.getDatabase()

	rows, err := db.Query(sqlText)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var commands Commands
	for rows.Next() {
		var command Command
		err := rows.Scan(&command.Name, &command.Response)
		if err != nil {
			log.Fatal(err)
		}
		commands = append(commands, command)
	}

	return commands
}

func (manager *SqLiteDataManager) addCommand(username, commandName, responseText string, cooldownSecond int) {
	db := manager.getDatabase()

	// Get user id from username.
	row := db.QueryRow("SELECT id FROM Users where username = $1", username)
	var userID int
	row.Scan(&userID)
	// First craete response
	result, _ := db.Exec("INSERT INTO Responses (name) VALUES ($1)", responseText)
	responseID, _ := result.LastInsertId()

	// Create command
	result, _ = db.Exec(
		"INSERT INTO Commands (name, user_id, response_id, cooldown_second) VALUES ($1, $2, $3, $4)",
		commandName, userID, responseID, cooldownSecond)
}
