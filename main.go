package jamie_chicken

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/scylladb/go-set/strset"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var api = slack.New(os.Getenv("SLACK_TOKEN"))

const ravesUserId = "U5Y1XU9UL"

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
		if IsGritty(strings.ToLower(ev.Text)) {
			_ = api.AddReaction("party-grit", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
		}
		if meetRavesGoal(ev) {
			_ = api.AddReaction("raves-goal", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
		}
		if words.Has("connect") {
			_ = api.AddReaction("dumpsterfire", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
		}
		if words.HasAny("mil", "mother", "blue", "eyes", "punnett") {
			_ = api.AddReaction("blue-eyes", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
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

func meetRavesGoal(m *slackevents.MessageEvent) bool {
	return m.User == ravesUserId && probablyNotGonnaHappenButMaybeCrazyJamieWillDoIt()
}

func probablyNotGonnaHappenButMaybeCrazyJamieWillDoIt() bool {
	return rand.Intn(100) < 1
}

func roll() bool {
	return rand.Intn(20) < 1
}
