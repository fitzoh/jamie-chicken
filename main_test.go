package jamie_chicken

import (
	"reflect"
	"testing"

	"github.com/scylladb/go-set/strset"
	"github.com/slack-go/slack/slackevents"
)

func Test_messageWords(t *testing.T) {
	message := &slackevents.MessageEvent{
		Text: "This is-a bunch of word's.",
	}
	want := strset.New("this", "is", "a", "bunch", "of", "word", "s")
	if got := messageWords(message); !reflect.DeepEqual(got, want) {
		t.Errorf("messageWords() = %v, want %v", got, want)
	}
}

func Test_hasAChickenWord(t *testing.T) {

	if hasAChickenWord(strset.New("dog")) {
		t.Errorf("dog is not a chicken")
	}

	if !hasAChickenWord(strset.New("chicken")) {
		t.Errorf("chicken is a chicken")
	}
}
