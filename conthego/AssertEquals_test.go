package conthego

import (
	"testing"
)

func Test(t *testing.T) {
	RunSpec()
}

func (f Fixture) GetGreeting() string {
	return "Hello World!"
}

func (f Fixture) GetPersonalisedGreeting(name string) string {
	return "Hello " + name + "!"
}

func (f Fixture) GetMultipleGreeting(name1 string, name2 string) string {
	return "Hello " + name1 + " and " + name2 + "!"
}
