package jamie_chicken

import (
	"fmt"
	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
	"github.com/scylladb/go-set/strset"
)

var model *text.NaiveBayes

//Jamie Dicken  1 day ago
//Oh yes, every time someone says sucks, awful, dumb, or fun, that will be me
var grittyWords = strset.New("grit", "sucks", "awful", "dumb", "fun", "stupid", "coldfusion", "outcomes", "jamie")
var regularWords = strset.New("awesome", "and", "some", "other", "regular", "words", "to", "balance", "the", "model")

func init() {
	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	model = text.NewNaiveBayes(stream, 2, base.OnlyWordsAndNumbers)
	go model.OnlineLearn(errors)

	grittyWords.Each(func(word string) bool {
		stream <- base.TextDatapoint{X: word, Y: 1}
		return true
	})
	regularWords.Each(func(word string) bool {
		stream <- base.TextDatapoint{X: word, Y: 0}
		return true
	})

	close(stream)

	for {
		err, _ := <-errors
		if err != nil {
			fmt.Printf("Error passed: %v", err)
		} else {
			// training is done!
			break
		}
	}
}

func IsGritty(text string) bool {
	return model.Predict(text) == 1
}
