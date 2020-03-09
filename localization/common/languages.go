package common

import "time"

// LocaleConfig language config
type LocaleConfig struct {
	LocaleID     string // BCP 47 representation of the language. All base languages should have one
	BaseLocaleID string // If empty, this language does not have fallback in case of missing translation
	DisplayName  string // Human-friendly display name for this language

	DateTimeToStringFunc func(time.Time) string // Function to convert Time object to string
	DurationToStringFunc func(Duration) string  // Function to convert duration to string

	InstallerLocale  *InstallerLocaleConfig
	BotCommandLocale *BotCommandLocaleConfig
}

func (locale *LocaleConfig) DateTimeToString(datetime time.Time) string {
	return locale.DateTimeToStringFunc(datetime)
}

func (locale *LocaleConfig) DurationToString(duration time.Duration) string {
	totalSeconds := int(duration.Seconds())
	d := NewDuration(totalSeconds)
	return locale.DurationToStringFunc(d)
}
