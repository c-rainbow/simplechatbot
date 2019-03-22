// Data Access Object to communicate with DB.

package models

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
/*type User struct {
	// Twitch ID
	TwitchID int64 `dynamo:"ID,hash"`
	// Twitch Username
	Username string `index:"Username-index"`
	// Twitch Oauth Token
	OauthToken string
}*/
