package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestAssertTrue(t *testing.T) {
	conthego.RunSpec(conthego.NewFixture(t, &FixtureAssertTrue{}))
}

type FixtureAssertTrue struct {
}

func (f FixtureAssertTrue) HowAreYou() bool {
	return true
}
