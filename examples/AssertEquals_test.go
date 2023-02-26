package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestAssertEquals(t *testing.T) {
	conthego.RunSpec(t, &AssertEqualsFixture{})
}

type AssertEqualsFixture struct {
}

func (f *AssertEqualsFixture) GetGreeting() string {
	return "Hello World! break test"
}

func (f *AssertEqualsFixture) GetPersonalisedGreeting(name string) string {
	return "Hello " + name + "!"
}

func (f *AssertEqualsFixture) GetMultipleGreeting(name1 string, name2 string) string {
	return "Hello " + name1 + " and " + name2 + "!"
}
