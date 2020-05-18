package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"strings"
	"testing"
)

func TestExecuteRows(t *testing.T) {
	conthego.RunSpec(t, &FixtureExecuteRows{})
}

type FixtureExecuteRows struct {
}

func (f FixtureExecuteRows) Split(name string) Name {
	split := strings.Split(name, " ")
	return Name{split[0], split[1]}
}
