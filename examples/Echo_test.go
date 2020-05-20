package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestEcho(t *testing.T) {
	conthego.RunSpec(t, &FixtureEcho{})
}

type FixtureEcho struct {
}

func (f *FixtureEcho) GetGreeting() string {
	return "Hello World!"
}

func (f *FixtureEcho) GetPersonalisedGreeting(name string) string {
	return "Hello " + name + "!"
}
