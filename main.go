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

func Handler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	event, e := slackevents.ParseEvent(bytes, slackevents.OptionNoVerifyToken())
	if e != nil {
		fmt.Println("error: ", e)
	}
	if event.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(bytes, &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	if event.Type == slackevents.CallbackEvent {
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			if roll() || ev.ChannelType == "im" {
				_ = api.AddReaction("chicken", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
			}
		case *slackevents.AppMentionEvent:
			_ = api.AddReaction("chicken", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
		}
	}
}

func roll() bool {
	return rand.Intn(20) < 1
}
