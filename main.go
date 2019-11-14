package jamie_chicken

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/scylladb/go-set/strset"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"unicode"
)

var api = slack.New(os.Getenv("SLACK_TOKEN"))

var chickenWords = strset.New("chicken", "chickens", "hen", "hens", "rooster", "roosters", "egg", "eggs", "bawk")

func main() {
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//Google cloud functions entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	event, e := slackevents.ParseEvent(bytes, slackevents.OptionNoVerifyToken())
	if e != nil {
		fmt.Println("error: ", e)
	}
	switch event.Type {
	case slackevents.URLVerification:
		verifyUrl(w, bytes)
	case slackevents.CallbackEvent:
		handleEvent(event.InnerEvent)
	}
}

func handleEvent(ie slackevents.EventsAPIInnerEvent) {
	switch ev := ie.Data.(type) {
	case *slackevents.MessageEvent:
		words := messageWords(ev)
		if roll() || ev.ChannelType == "im" || hasAChickenWord(words) {
			_ = api.AddReaction("chicken", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
		}
	case *slackevents.AppMentionEvent:
		_ = api.AddReaction("chicken", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
	}
}

//Handle an initial slack verification request... Should only happen on install
func verifyUrl(w http.ResponseWriter, bytes []byte) {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal(bytes, &r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text")
	_, _ = w.Write([]byte(r.Challenge))
}

func messageWords(m *slackevents.MessageEvent) *strset.Set {
	lowerMessage := strings.ToLower(m.Text)
	lowerWords := strings.FieldsFunc(lowerMessage, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	return strset.New(lowerWords...)
}

func hasAChickenWord(w *strset.Set) bool {
	return !strset.Intersection(chickenWords, w).IsEmpty()
}

func roll() bool {
	return rand.Intn(20) < 1
}
