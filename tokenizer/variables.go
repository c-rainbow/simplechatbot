package tokenizer

// Variables that can be used
type Variable struct {
	Name        string
	Type        string
	Description string
	Enabled     bool // TODO: Is it possible to enable in certain channels?
}

// Variable types
// TODO: This can be int, I think
const (
	ChatType               = "chat"               // Values can be found in the chat message
	StreamAPIType          = "streamAPI"          // Twitch Stream API call
	UserAPIType            = "userAPI"            // Twitch Users API call
	SimpleType             = "simple"             // No external API call
	OverwatchAPIType       = "overwatchAPI"       // Overwatch API call
	LeagueOfLegendsAPIType = "LeagueOfLegendsAPI" // Riot LoL API call
)

/*
$user: 채팅 친 사람 display name
$user_id: 채팅 친 사람 twitch username (not display name)
$display_name: 채팅 친 사람 display name
$channel: 채팅이 올라온 채널 이름
$subscribe: 구독 링크 (http://twitch.tv/subs/$channel)
$commands: 명령어 페이지 링크
*/
var (
	User = Variable{
		Name:        "user",
		Type:        ChatType,
		Description: "Chatter's display name",
		Enabled:     true,
	}
	UserID = Variable{
		Name:        "user_id",
		Type:        ChatType,
		Description: "Chatter's Twitch user ID (integer)",
		Enabled:     true,
	}
	DisplayName = Variable{
		Name:        "display_name",
		Type:        ChatType,
		Description: "Chatter's Twitch display name", // Currently same as $user
		Enabled:     true,
	}
	Channel = Variable{
		Name:        "channel",
		Type:        ChatType,
		Description: "The channel's name where the chat was posted",
		Enabled:     true,
	}
	// TODO: Automatically enable when the streamer is affiliate/partner?
	SubscribeLink = Variable{
		Name:        "subscription_link",
		Type:        ChatType,
		Description: "The channel's subscription link http://twitch.tv/subs/$(channel)",
		Enabled:     true,
	}
	// TODO: Find a good domain
	Commands = Variable{
		Name:        "commands",
		Type:        ChatType,
		Description: "Link to the commands webpage http://....../$(channel)",
		Enabled:     false,
	}
)

/*
$title: 채널 방 제목
$game: 현재 방송중인 게임 이름
$uptime: 업타임
$uptime_at: 방송 시작 시간
$viewer_count: 시청자 수
$follower_count: 팔로워 수
$subscriber_count: 구독자 수
$status: 채널 상태 (방송중, 오프라인, etc)
*/
var (
	Title = Variable{
		Name:        "title",
		Type:        StreamAPIType,
		Description: "The Channel's title",
		Enabled:     false,
	}
	Game = Variable{
		Name:        "game",
		Type:        StreamAPIType,
		Description: "The game being played in the stream",
		Enabled:     false,
	}
	Uptime = Variable{
		Name:        "uptime",
		Type:        StreamAPIType,
		Description: "Stream uptime",
		Enabled:     false,
	}
	UptimeAt = Variable{
		Name:        "uptime_at",
		Type:        StreamAPIType,
		Description: "Stream start time",
		Enabled:     false,
	}
	ViewerCount = Variable{
		Name:        "viewer_count",
		Type:        StreamAPIType,
		Description: "Number of current viewers",
		Enabled:     false,
	}
	FollowerCount = Variable{
		Name:        "follower_count",
		Type:        StreamAPIType,
		Description: "Number of followers",
		Enabled:     false,
	}
	// TODO: This variable does not use the same API as other ones
	SubscriberCount = Variable{
		Name:        "subscriber_count",
		Type:        StreamAPIType,
		Description: "Number of subscribers",
		Enabled:     false,
	}
)

/*
$follow_age (username) : 사용자가 팔로우한 기간
$follow_start_date (username): username이 팔로우한 시각
*/
var (
	FollowAge = Variable{
		Name:        "follow_age",
		Type:        UserAPIType,
		Description: "How long the user has followed the channel",
		Enabled:     false,
	}
	FollowStartDate = Variable{
		Name:        "follow_start_date",
		Type:        UserAPIType,
		Description: "When the user followed the channel",
		Enabled:     false,
	}
)

/*
$current_time (location): location의 현재 시각
$countdown (datetime): datetime까지 남은 시간
$countup (datetime): datetime으로부터 지난 시간
$count (variable, inc=1): call이 될 때마다 variable을 inc씩 increment
$rand (start, end): start 와 end 사이의 (둘다 integer, 양쪽으로 inclusive) 랜덤한 숫자 리턴
*/
var (
	// TODO: Which arguments are needed to use this variable?
	// For example, should user provide only the location (city/town name), or
	// time zone in different formats (e.g. UTC-01, EST, America/Los_Angeles)
	// Also, how to get the current time correctly in case of daylight saving?
	CurrentTime = Variable{
		Name:        "time",
		Type:        SimpleType, // Not simple!
		Description: "Current time of the location (location argument needed)",
		Enabled:     false,
	}
	Countdown = Variable{
		Name:        "countdown",
		Type:        SimpleType,
		Description: "Countdown to the provided datetime (argument needed)",
		Enabled:     false,
	}
	Countup = Variable{
		Name:        "countup",
		Type:        SimpleType,
		Description: "Count up to the provided datetime (argument needed)",
		Enabled:     false,
	}
	Count = Variable{
		Name:        "count",
		Type:        SimpleType,
		Description: "Count how many the argument name is called (argument needed)",
		Enabled:     false,
	}
	Rand = Variable{
		Name:        "rand",
		Type:        SimpleType,
		Description: "Return random int between [a, b] inclusive (arguments needed)",
		Enabled:     false,
	}
)

/*
$overwatch (battletag): battletag의 오버워치 현 점수
$leagueoflegends (summoner_name, region): LOL 현 티어
*/
var (
	Overwatch = Variable{
		Name:        "overwatch",
		Type:        OverwatchAPIType,
		Description: "Current Overwatch rank in competitive play",
		Enabled:     false,
	}
	LeagueOfLegends = Variable{
		Name:        "league_of_legends",
		Type:        LeagueOfLegendsAPIType,
		Description: "Current League of Legends tier",
		Enabled:     false,
	}
)
