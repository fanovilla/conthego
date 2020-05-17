package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestReturnStruct(t *testing.T) {
	conthego.RunSpec(conthego.NewFixture(&FixtureReturnStruct{}))
}

type FixtureReturnStruct struct {
}

type Name struct {
	First string
	Last  string
}

func (f FixtureReturnStruct) BuildName() Name {
	return Name{"Ryan", "Liam"}
}
