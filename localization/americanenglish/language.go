package americanenglish

import (
	"strings"

	"github.com/c-rainbow/simplechatbot/localization/common"
	"golang.org/x/text/language"
)

var (
	// LocaleID BCP 47 representation of the language
	LocaleID = strings.ToLower(language.AmericanEnglish.String())

	// Config American English locale config
	Config = common.LocaleConfig{
		LocaleID:    LocaleID,
		DisplayName: "American English",

		DateTimeToStringFunc: DateTimeToString,
		DurationToStringFunc: DurationToString,

		InstallerLocale:  &installerLocale,
		BotCommandLocale: &botCommandLocale,
	}
)
