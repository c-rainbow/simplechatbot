package americanenglish

import "github.com/c-rainbow/simplechatbot/localization/common"

var (
	installerLocale = common.InstallerLocaleConfig{
		Messages: common.InstallerMessages{
			NoHelixClient:                "트위치 Helix API에 연결할 수 없습니다. Client ID를 확인해 주세요",
			NoKrakenClient:               "트위치 Kraken API에 연결할 수 없습니다. Client ID를 확인해 주세요",
			TwitchUsersAPIError:          "트위치 Helix Users API에 연결할 수 없습니다",
			TwitchBotAccountNotFound:     "봇 트위치 계정을 찾을 수 없습니다",
			TwitchChannelAccountNotFound: "채널을 찾을 수 없습니다",
			ChatServerSuccessfulLogin:    "채팅 서버에 연결을 실패하였습니다",
			ChatServerFailedLogin:        "채팅창 로그인에 실패하였습니다. OAuth 토큰을 확인해 주세요",
			DynamoDBConnectionError:      "DynamoDB에 연결할 수 없습니다",
		},
	}
)
