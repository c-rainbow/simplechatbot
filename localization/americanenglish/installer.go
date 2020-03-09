package americanenglish

import "github.com/c-rainbow/simplechatbot/localization/common"

var (
	installerLocale = common.InstallerLocaleConfig{
		Messages: common.InstallerMessages{
			NoHelixClient:                "Cannot connect to Twitch Helix API. Please check ClientID",
			NoKrakenClient:               "Cannot connect to Twitch Kraken API. Please check ClientID",
			TwitchUsersAPIError:          "Cannot connect to Twitch Helix User API",
			TwitchBotAccountNotFound:     "Cannot find the bot's Twitch account",
			TwitchChannelAccountNotFound: "Cannot find the streamer's Twitch channel",
			ChatServerSuccessfulLogin:    "Cannot connect to Twitch chat server",
			ChatServerFailedLogin:        "Login to chat failed. Please check your OAuth token",
			DynamoDBConnectionError:      "Cannot connect to DynamoDB",
		},
	}
)
