package jamie_chicken

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type emotionResponse struct {
	EmotionScore emotionScore `json:"EmotionScore"`
	Emotion      string       `json:"Emotion"`
}

type emotionScore struct {
	Anger          float32 `json:"Anger"`
	Neutral        float32 `json:"Neutral"`
	Others         float32 `json:"Others"`
	Happy          float32 `json:"Happy"`
	Joy            float32 `json:"Joy"`
	Affection      float32 `json:"Affection"`
	Sad            float32 `json:"Sad"`
	Disappointment float32 `json:"Disappointment"`
}

type conversation struct {
	Conversation []string `json:"conversation"`
	Lang         string   `json:"lang"`
	IgnoreFirst  bool     `json:"ignore_first,omitempty"`
}

type sentimClient struct {
	auth_token string
	base_url   string
	logger     *logrus.Logger
}

func NewSentimClient() *sentimClient {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	return &sentimClient{
		base_url: "https://api.sentimllc.com/",
		logger:   l,
	}
}

func (c *sentimClient) IsGritty(message string) bool {
	return c.isEmotion("Anger", message)
}

func (c *sentimClient) IsDisappointed(message string) bool {
	return c.isEmotion("Disappointment", message)
}

func (c *sentimClient) isEmotion(emotion string, message string) bool {
	emot, err := c.getEmotion(message)
	if err != nil {
		c.logger.Debugf("something went wrong with getting emotion from the api ¯\\_(ツ)_/¯")
		c.logger.Errorf("I supposed you want to see the error: %s", err)
		return false
	}

	c.logger.Debugf("emotion: %s", emot.Emotion)
	return strings.ToUpper(emot.Emotion) == strings.ToUpper(emotion)
}

func (c *sentimClient) getEmotion(message string) (*emotionResponse, error) {
	conversation_req := conversation{
		Lang:         "en",
		Conversation: []string{message},
	}
	bytez, err := json.Marshal(conversation_req)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/emotion/single", c.base_url), bytes.NewBuffer(bytez))
	if err != nil {
		return nil, err
	}
	c.authIfNecessary()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.auth_token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status `%d` for emotion api call", resp.StatusCode)
	}

	var emot_resp *emotionResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&emot_resp); err != nil {
		return nil, err
	}
	return emot_resp, nil
}

func (c *sentimClient) authIfNecessary() {
	if c.auth_token == "" {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/token", c.base_url), nil)
		if err != nil {
			c.logger.Errorf("error building request: %s", err)
			return
		}

		params := url.Values{}
		params.Add("client_id", os.Getenv("SENTIM_CLIENTID"))
		params.Add("client_secret", os.Getenv("SENTIM_CLIENTSECRET"))
		req.URL.RawQuery = params.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.logger.Errorf("post to api failed for some reason: %s", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			c.logger.Errorf("did not receive 200, got `%d` (aka check creds)", resp.StatusCode)
			return
		}
		defer resp.Body.Close()

		var tok string
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&tok); err != nil {
			c.logger.Errorf("decoding response didn't work: %s", err)
			return
		}
		c.auth_token = tok
	}
}
