package install

import (
	"errors"
	"log"

	"github.com/c-rainbow/simplechatbot/flags"

	"github.com/c-rainbow/simplechatbot/api/helix"
	"github.com/c-rainbow/simplechatbot/client"
	"github.com/go-ini/ini"
)

var (
	ErrTwitchAccoutNotFound = errors.New("Twitch Account Not Found")
)

/*

각각의 단계에서 안되면 실패

1. 설정파일 읽어오고 확인하기. 필요한 값들 모두 있는지 최소한 형식 (int, string, etc) 은 맞는지.

2. 설정파일 봇 존재하는 계정인지 확인
3. 설정파일 봇으로 채팅서버에 로그인 해보기. 테스트 채널에서 채팅 쳐보기 (어떻게 잘 되었는지 확인?)

4. 설정파일 채널 존재하는 계정인지 확인

5. DynamoDB 접속해 보기
6. DB에 Bots 테이블 만들기
7. DB에 Channels 테이블 만들기
8. 기본 봇 데이타 넣기
9. 기본 채널 데이타 넣기
10. 채널에 봇 넣기

*/

var (
	InstallerMessagesKorean = &InstallerMessages{
		NoHelixClient:                "트위치 API에 연결할 수 없습니다",
		TwitchUsersAPIError:          "트위치 Users API 에서 정보를 가져올 수 없습니다",
		TwitchBotAccountNotFound:     "봇 계정을 찾을 수 없습니다",
		TwitchChannelAccountNotFound: "채널 계정을 찾을 수 없습니다",
		ChatServerSuccessfulLogin:    "채팅 서버에 성공적으로 접속하였습니다. 3초 후 접속을 종료합니다",
		ChatServerFailedLogin:        "채팅 서버 접속에 실패하였습니다. 설정을 확인해 주세요.",
		DynamoDBConnectionError:      "DynamoDB에 접속 중 오류가 발생했습니다",
	}
)

func NewInstallerKor() *Installer {
	iniFile, err := ini.Load(flags.InstallationConfigFile)
	if err != nil {
		log.Fatalln("설정 파일을 읽어오는데 실패하였습니다.", err.Error())
	}

	installer := &Installer{
		iniFile: iniFile, ircClientFunc: client.NewTwitchClient, helixClientFunc: helix.NewHelixClient,
		messages: InstallerMessagesKorean}
	return installer
}
