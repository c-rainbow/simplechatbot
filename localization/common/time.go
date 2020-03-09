package common

// Duration duration struct needed for displaying stream uptime, follow duration, etc
type Duration struct {
	Day    int
	Hour   int
	Minute int
	Second int
}

// IsZero returns true if all fields are zero.
func (duration *Duration) IsZero() bool {
	return duration.Day == 0 && duration.Hour == 0 && duration.Minute == 0 && duration.Second == 0
}

// NewDuration converts total seconds to Duration object.
func NewDuration(totalSeconds int) Duration {
	duration := Duration{}
	duration.Second = totalSeconds % 60
	totalSeconds /= 60
	duration.Minute = totalSeconds % 60
	totalSeconds /= 60
	duration.Hour = totalSeconds % 24
	totalSeconds /= 24
	duration.Day = totalSeconds

	return duration
}
