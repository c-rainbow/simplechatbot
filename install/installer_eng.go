package install

import (
	"log"

	"github.com/c-rainbow/simplechatbot/api/helix"
	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/flags"
	"github.com/go-ini/ini"
)

var (
	InstallerMessagesEnglish = &InstallerMessages{
		NoHelixClient:                "Cannot connect to Twitch API.",
		TwitchUsersAPIError:          "Error getting data from Twitch Users API.",
		TwitchBotAccountNotFound:     "Twitch Bot account not found.",
		TwitchChannelAccountNotFound: "Twitch account for the channel not found.",
		ChatServerSuccessfulLogin:    "Successfully logged in to the chat server. Disconnecting..",
		ChatServerFailedLogin:        "Failed to log in to the chat server. Please make sure the config is correct.",
		DynamoDBConnectionError:      "Error while connecting to DynamoDB.",
	}
)

func NewInstallerEng() *Installer {
	iniFile, err := ini.Load(flags.InstallationConfigFile)
	if err != nil {
		log.Fatalln("Failed to load config file.", err.Error())
	}

	installer := &Installer{
		iniFile: iniFile, ircClientFunc: client.NewTwitchClient, helixClientFunc: helix.NewHelixClient,
		messages: InstallerMessagesEnglish}
	return installer
}
