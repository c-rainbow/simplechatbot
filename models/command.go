package models

// User levels for command.
// Externally, only Everyone, Moderator, Streamer will be available.
const (
	// Everyone can call the command
	PermissionEveryone int = 1 << 0
	// Follower only
	// PermissionFollower int = 1 << 1
	// Subscriber only
	PermissionSubscriber int = 1 << 2
	// Vip only
	// PermissionVip int = 1 << 3
	// Moderator only
	PermissionModerator int = 1 << 4
	// Streamer only
	PermissionStreamer int = 1 << 5
)

var (
	// Default response key in resopnse map
	DefaultResponseKey        = "DEFAULT"
	DefaultFailureResponseKey = "FAILURE"
)

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
	// Chat plugin type: response, add_command, edit_command, delete_command, etc
	PluginType string
	// Bot's Response, in parsed form
	Responses map[string]ParsedResponse // `dynamo:",set"`
	// Cooldown in seconds
	CooldownSecond int
	// Permission bitset as integer
	Permission int
	// True if enabled
	Enabled bool
	// Group of commands. Commands can be active/inactive per group.
	Group string
}

// Create a new Command object with inputs.
func NewCommand(botID int64, channelID int64, name string, plugintype string, defaultResponse *ParsedResponse,
	failureReaponse *ParsedResponse) *Command {
	// Build minimal response map
	responseMap := map[string]ParsedResponse{
		DefaultResponseKey:        *defaultResponse,
		DefaultFailureResponseKey: *failureReaponse,
	}
	// Some properties of Command are initialized.
	command := Command{
		BotID:          botID,
		ChannelID:      channelID,
		Name:           name,
		PluginType:     plugintype,
		Responses:      responseMap,
		Permission:     PermissionEveryone,
		CooldownSecond: 5,
		Enabled:        true,
		Group:          "",
	}
	return &command
}
