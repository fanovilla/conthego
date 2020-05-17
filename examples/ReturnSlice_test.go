package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestReturnSlice(t *testing.T) {
	conthego.RunSpec(t, &FixtureReturnSlice{})
}

var list *[]string

type FixtureReturnSlice struct {
}

func (f FixtureReturnSlice) BuildList() []string {
	list = &[]string{"Carl", "Ryan", "Liam"}
	return *list
}

func (f FixtureReturnSlice) GetFirstEntry() string {
	return (*list)[0]
}
