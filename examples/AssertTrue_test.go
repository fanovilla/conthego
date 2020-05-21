package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestAssertTrue(t *testing.T) {
	conthego.RunSpec(t, &AssertTrueFixture{})
}

type AssertTrueFixture struct {
}

func (f *AssertTrueFixture) IsTheWorldRound() bool {
	return true
}

func (f *AssertTrueFixture) IsTheWorldFlat() bool {
	return false
}
