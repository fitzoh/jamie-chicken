package jamie_chicken_test

import (
	"testing"

	"github.com/fitzoh/jamie_chicken"
	"github.com/stretchr/testify/assert"
)

func TestItDoesAThing(t *testing.T) {
	c := jamie_chicken.NewSentimClient()
	assert.Equal(t, c.IsGritty("wtf, I actually kind of know go now :party-dinosaur:"), true)
}

func TestHowAngryDoIHaveToGetToBeGritty(t *testing.T) {
	c := jamie_chicken.NewSentimClient()
	assert.Equal(t, c.IsGritty("the fuck do I have to say to be angry"), true)
}

func TestDisappointmentIsAThing(t *testing.T) {
	c := jamie_chicken.NewSentimClient()
	assert.Equal(t, c.IsDisappointed("what do you mean? the outcomes I know gives wide latitude and lets teams work independently to drive produc….. hahaha couldn’t even finish the bullshit"), false)
}
