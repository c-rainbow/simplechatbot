// Data Access Object to communicate with DB.

package models

// User levels for command.
// Externally, only Everyone, Moderator, Stremer will be available.
const (
	// Everyone can call the command
	Everyone int = 1 << 0
	// Follower only
	Follower int = 1 << 1
	// Subscriber only. VIP is also considered to be subscriber.
	Subscriber int = 1 << 2
	// Vip only
	Vip int = 1 << 3
	// Moderator only
	Moderator int = 1 << 4
	// Streamer only
	Streamer int = 1 << 5
)

// Bot contains info about Twitch chatbots.
type Bot struct {
	// Bot's Twitch ID
	TwitchID int64 `dynamo:"ID,hash"`
	// Bot's Twitch username
	Username string `index:"Username-index,hash"`
	// Twitch Oauth token
	OauthToken string
}

// Channel describes Twitch channel
type Channel struct {
	// Channel's Twitch ID
	TwitchID int64 `dynamo:"ID,hash"`
	// Channel's Twitch Username
	Username string `index:"Username-index,hash"`
	// Channel's Twitch display name
	DisplayName string
	// Bots which join this channel
	BotIDs []int64 `dynamo:",set"`
	// Commands of the channel
	Commands map[string]Command // `dynamo:",set"`
}

// User describes Twitch user. Not sure when this will be used.
type User struct {
	// Twitch ID
	TwitchID int64 `dynamo:"ID,hash"`
	// Twitch Username
	Username string `index:"Username-index"`
	// Twitch Oauth Token
	OauthToken string
}

// Command chatbot commands
type Command struct {
	// Generated (not by DB) unique ID for command
	// UUID string `dynamo:"ID,hash"`
	// Bot's Twitch ID
	BotID int64
	// Channel's Twitch ID
	ChannelID int64
	// Command name
	Name string
	// Bot's Response, in parametrized string
	Responses map[string]string `dynamo:",set"`
	// Cooldown in seconds
	CooldownSecond int
	// True if enabled
	Enabled bool
	// Group of commands. Commands can be active/inactive per group.
	Group string
}
