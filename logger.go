package jamie_chicken

import (
	"fmt"

	"github.com/slack-go/slack"
)

type Logger struct {
	slack   *slack.Client
	channel string
}

func NewLogger(s *slack.Client, ch string) *Logger {
	return &Logger{
		slack:   s,
		channel: ch,
	}
}

func (l *Logger) Error(msg string) {
	_, _, err := l.slack.PostMessage(
		l.channel,
		slack.MsgOptionText(fmt.Sprintf("```%s```", msg), false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Printf("uhhhhh...... now what")
		fmt.Printf(err.Error())
	}
}
