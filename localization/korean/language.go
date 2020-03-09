package korean

import (
	"strings"

	"github.com/c-rainbow/simplechatbot/localization/common"
	"golang.org/x/text/language"
)

var (
	// LocaleID BCP 47 representation of the language
	LocaleID = strings.ToLower(language.Korean.String())

	// Config Korean locale config
	Config = common.LocaleConfig{
		LocaleID:    LocaleID,
		DisplayName: "한국어",

		DateTimeToStringFunc: DateTimeToString,
		DurationToStringFunc: DurationToString,

		InstallerLocale:  &installerLocale,
		BotCommandLocale: &botCommandLocale,
	}
)
