package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestEcho(t *testing.T) {
	conthego.RunSpec(t, &SetAndEchoFixture{})
}

type SetAndEchoFixture struct {
}

func (f *SetAndEchoFixture) GetGreeting() string {
	return "Hello World!"
}

func (f *SetAndEchoFixture) GetPersonalisedGreeting(name string) string {
	return "Hello " + name + "!"
}
