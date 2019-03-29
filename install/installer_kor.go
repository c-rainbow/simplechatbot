package install

type InstallerKor struct {
	configFile string
}

func InstallKor() {

}

func ReadConfig() error {

	return nil
}

func TryAccessingDynamoDB() error {
	return nil
}

func CheckIfChannelAccountExists() error {
	return nil
}

func CheckIfBotAccountExists() error {
	return nil
}

func AddDefaultBot() error {
	return nil
}
func AddDefaultChannel() error {
	return nil
}
func AddBotToChannel() error {
	return nil
}

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
