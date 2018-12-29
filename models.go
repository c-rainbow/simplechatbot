// Models to communicate with DB.

package simplechatbot

// User contains Twitch user info.
type User struct {
	// User's Twitch username
	username string `json:"username"`
}

// Users multiple users.
type Users []User

// Bot contains info about Twitch chatbots. Usually there is only one.
type Bot struct {
	// Bot's Twitch username
	name string `json:"name"`
	// Twitch Oauth token
	oauthToken string `json:"oauthToken`
}

// Bots multiple bots.
type Bots []Bot

// Command contains info about chatbot commands.
type Command struct {
	// Bot command name
	name string `json:"name"`
	// Bot's response, in parametrized string
	response string `json:"response"`
	// Cooldown in second
	cooldownSecond int `json:"cooldownSecond"`
}

// Commands multiple commands.
type Commands []Command

// UserCommand which user has which commands
type UserCommand struct {
	user     User
	commands Commands
}

// UserCommands multiple UserCommand
type UserCommands []UserCommand
