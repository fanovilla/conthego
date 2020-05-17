package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestExecute(t *testing.T) {
	conthego.RunSpec(conthego.NewFixture(&FixtureExecute{}))
}

type FixtureExecute struct {
}

func (f FixtureExecute) GetPersonalisedGreeting(name string) string {
	return "Hello " + name + "!"
}
