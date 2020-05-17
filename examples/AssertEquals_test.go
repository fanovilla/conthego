package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestAssertEquals(t *testing.T) {
	conthego.RunSpec(t, &FixtureAssertEquals{})
}

type FixtureAssertEquals struct {
}

func (f FixtureAssertEquals) GetGreeting() string {
	return "Hello World!"
}

func (f FixtureAssertEquals) GetPersonalisedGreeting(name string) string {
	return "Hello " + name + "!"
}

func (f FixtureAssertEquals) GetMultipleGreeting(name1 string, name2 string) string {
	return "Hello " + name1 + " and " + name2 + "!"
}
