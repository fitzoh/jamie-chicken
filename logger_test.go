package jamie_chicken_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fitzoh/jamie_chicken"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
)

func TestLoggerErrorJustPostsMessageIGuess(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	api := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))

	ch := "super-secret-jamie-chicken-logging-channel"
	l := jamie_chicken.NewLogger(api, ch)

	msg := "I have no idea what I'm doing"
	l.Error(msg)

	out := s.GetSeenOutboundMessages()
	assert.Equal(t, 1, len(out))
	var m = slack.Message{}
	json.Unmarshal([]byte(out[0]), &m)
	assert.Equal(t, fmt.Sprintf("```%s```", msg), m.Text)
	assert.Equal(t, ch, m.Channel)
}
