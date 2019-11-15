package jamie_chicken

import (
    "testing"
)

func TestFitzIsNotSoGritty(t *testing.T) {
    message := "fitz is awesome"

    if IsGritty(message) {
        t.Errorf("expected non-gritty sentiment")
    }
}

func TestRandomTextIsNotGritty(t *testing.T) {
    message := "some random text"

    if IsGritty(message) {
        t.Errorf("expected non-gritty sentiment")
    }
}

func TestColdAndFusionButNotColdFusion(t *testing.T) {
    message := "it's really cold out, we need to figure out nuclear fusion soon"

    if IsGritty(message) {
        t.Errorf("expected non-gritty sentiment")
    }
}

func TestColdfusionIsGritty(t *testing.T) {
    message := "coldfusion fucking sucks"

    if !IsGritty(message) {
        t.Errorf("expected gritty sentiment")
    }
}

func TestSucksDumbStupidIsGritty(t *testing.T) {
    message := "sucks this dumb message is about a stupid chicken or whatever"

    if !IsGritty(message) {
        t.Errorf("expected gritty sentiment")
    }
}
