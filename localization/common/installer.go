package common

// InstallerLocaleConfig locale config for installer
type InstallerLocaleConfig struct {
	Messages InstallerMessages
}

// InstallerMessages Installer messages
type InstallerMessages struct {
	NoHelixClient                string
	NoKrakenClient               string
	TwitchUsersAPIError          string
	TwitchBotAccountNotFound     string
	TwitchChannelAccountNotFound string
	ChatServerSuccessfulLogin    string
	ChatServerFailedLogin        string
	DynamoDBConnectionError      string
}
