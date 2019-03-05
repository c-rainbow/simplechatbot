// Models to communicate with DB.

package simplechatbot

// User contains Twitch user info.
type User struct {
	// User's Twitch Username
	Username string `json:"username"`
}

// Users multiple users.
type Users []User

// Bot contains info about Twitch chatbots. Usually there is only one.
type Bot struct {
	// Bot's Twitch username
	Name string `json:"name"`
	// Twitch Oauth token
	OauthToken string `json:"oauthToken`
}

// Bots multiple bots.
type Bots []Bot

// Command contains info about chatbot commands.
type Command struct {
	// Bot command name
	Name string `json:"name"`
	// Bot's Response, in parametrized string
	Response string `json:"response"`
	// Cooldown in second
	CooldownSecond int `json:"cooldownSecond"`
}

// Commands multiple commands.
type Commands []Command

// UserCommand which user has which commands
type UserCommand struct {
	User     User
	Commands Commands
}

// UserCommands multiple UserCommand
type UserCommands []UserCommand
