package korean

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/c-rainbow/simplechatbot/localization/common"
)

// DurationToString string representation of Duration in Korean
func DurationToString(duration common.Duration) string {
	// Shortcut for zero duration
	if duration.IsZero() {
		return "0초"
	}
	slices := make([]string, 0, 4)
	// Using strconv.Itoa is faster than fmt.Sprintf in this case
	if duration.Day != 0 {
		slices = append(slices, strconv.Itoa(duration.Day)+"일")
	}
	if duration.Hour != 0 {
		slices = append(slices, strconv.Itoa(duration.Hour)+"시간")
	}
	if duration.Minute != 0 {
		slices = append(slices, strconv.Itoa(duration.Minute)+"분")
	}
	if duration.Second != 0 {
		slices = append(slices, strconv.Itoa(duration.Second)+"초")
	}

	return strings.Join(slices, " ")
}

// DateTimeToString returns string representation of date time.
// Example of output: 2020년 3월 8일 9:03PM
func DateTimeToString(t time.Time) string {
	// US English AM/PM notation, as well as 오전/오후, seems to be commonly used in Korea
	// Use US English notation for simplicity of code
	localTime := t.Format(time.Kitchen)
	return fmt.Sprintf("%d년 %d월 %d일 %s", t.Year(), t.Month(), t.Day(), localTime)
}
