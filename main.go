package jamie_chicken

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
)

var api = slack.New(os.Getenv("SLACK_TOKEN"))

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

//We got an actual event of some kind, let's (maybe) do something stupid
func handleEvent(ie slackevents.EventsAPIInnerEvent) {
	switch ev := ie.Data.(type) {
	case *slackevents.MessageEvent:
		if roll() || ev.ChannelType == "im" {
			_ = api.AddReaction("chicken", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
		}
	case *slackevents.AppMentionEvent:
		_ = api.AddReaction("chicken", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
	}
}

func roll() bool {
	return rand.Intn(20) < 1
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
