package localization

import (
	"github.com/c-rainbow/simplechatbot/localization/americanenglish"
	"github.com/c-rainbow/simplechatbot/localization/common"
	"github.com/c-rainbow/simplechatbot/localization/korean"
)

/*

Localization package

There are two situations where localization is needed. A locale does not need to support both cases.

(1) Installer language

(2) Response text for bot commands
	(2-1) Default response text for fallback when response text does not exist for the given response key
	(2-2) Time-related localization when resolving variables like $(uptime), $(time), $(since), etc
	(2-3) [FUTURE] currency-related localization, etc?


The standard locales have BCP 47 locale ID, but users may add custom locales,
which is completely new or based on existing locale

TODO: Read locale configs from data files (JSON, YAML, etc)

*/

var (
	// SupportedLocales locales to support. All locales should be registered here
	SupportedLocales = []*common.LocaleConfig{
		&americanenglish.Config, // en-us
		&korean.Config,          // ko
	}

	localeMap = make(map[string]*common.LocaleConfig)
)

// GetAllInstallerLocales all supported installer locales
func GetAllInstallerLocales() []*common.LocaleConfig {
	locales := make([]*common.LocaleConfig, 0, 4)
	for _, locale := range SupportedLocales {
		if locale.InstallerLocale != nil {
			locales = append(locales, locale)
		}
	}
	return locales
}

// GetAllBotResponseLocales all supported bot response locales
func GetAllBotResponseLocales() []*common.LocaleConfig {
	locales := make([]*common.LocaleConfig, 0, 4)
	for _, locale := range SupportedLocales {
		if locale.BotCommandLocale != nil {
			locales = append(locales, locale)
		}
	}
	return locales
}

// GetBotResponseLocale bot response locale from LocaleID
func GetBotResponseLocale(localeID string) *common.LocaleConfig {
	// Not initialized yet
	if len(localeMap) == 0 {
		for _, config := range GetAllBotResponseLocales() {
			localeMap[config.LocaleID] = config
		}
	}
	// Returns nil if no locale
	return localeMap[localeID]
}
