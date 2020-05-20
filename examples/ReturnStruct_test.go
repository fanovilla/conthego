package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestReturnStruct(t *testing.T) {
	conthego.RunSpec(t, &FixtureReturnStruct{})
}

type FixtureReturnStruct struct {
}

func (f FixtureReturnStruct) BuildName() Name {
	return Name{"Ryan", "Liam"}
}
