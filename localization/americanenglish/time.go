package americanenglish

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/c-rainbow/simplechatbot/localization/common"
)

// DurationToString string representation of Duration in American English
func DurationToString(duration common.Duration) string {
	// Shortcut for zero duration
	if duration.IsZero() {
		return "0 seconds"
	}
	slices := make([]string, 0, 8)
	if duration.Day != 0 {
		slices = pluralized(slices, duration.Day, "day", "days")
	}
	if duration.Hour != 0 {
		slices = pluralized(slices, duration.Hour, "hour", "hours")
	}
	if duration.Minute != 0 {
		slices = pluralized(slices, duration.Minute, "minute", "minutes")
	}
	if duration.Second != 0 {
		slices = pluralized(slices, duration.Second, "second", "seconds")
	}

	return strings.Join(slices, " ")
}

// DateTimeToString returns string representation of date time.
// Go language does not have a pre-defined format for US-style representation
// Example of output: March 8, 2020 9:03PM
func DateTimeToString(t time.Time) string {
	monthName := t.Month().String()
	localTime := t.Format(time.Kitchen)
	return fmt.Sprintf("%s %d, %d %s", monthName, t.Day(), t.Year(), localTime)
}

// Append number and singular word if 1, plural word if number is not one.
func pluralized(slice []string, number int, singular string, plural string) []string {
	if number == 1 {
		return append(slice, strconv.Itoa(number), singular)
	}
	return append(slice, strconv.Itoa(number), plural)
}
