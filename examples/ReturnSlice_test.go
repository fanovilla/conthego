package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestReturnSlice(t *testing.T) {
	conthego.RunSpec(t, &FixtureReturnSlice{})
}

type FixtureReturnSlice struct {
	list *[]string
}

func (f *FixtureReturnSlice) BuildList() []string {
	f.list = &[]string{"Carl", "Ryan", "Liam"}
	return *f.list
}

func (f *FixtureReturnSlice) GetFirstEntry() string {
	return (*f.list)[0]
}
