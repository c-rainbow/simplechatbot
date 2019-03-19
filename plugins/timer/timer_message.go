package timer

// Package for plugin that sends timer messages to channels.
// Note that this is the manager, not each mesage.
// Usually there would be only one timer message manager per bot

type TimerMessageManager struct {
}

func (manager *TimerMessageManager) Run() {

}

// Read from repository of all
func (manager *TimerMessageManager) Read() {}

func (manager *TimerMessageManager) Action() {}

func (manager *TimerMessageManager) Output() {}

func (manager *TimerMessageManager) Store() {}

// TODO: not string, but ParsedTimer message
func (manager *TimerMessageManager) AddTimer(channel string, message string) {}

func (manager *TimerMessageManager) RemoveTimer(channel string, message string) {}
