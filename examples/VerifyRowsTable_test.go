package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"strings"
	"testing"
)

func TestVerifyRowsTable(t *testing.T) {
	conthego.RunSpec(t, &FixtureVerifyRowsTable{})
}

type FixtureVerifyRowsTable struct {
	users []string
}

func (f *FixtureVerifyRowsTable) SetUpUser(user string) []string {
	f.users = append(f.users, user)
	return f.users
}

func (f *FixtureVerifyRowsTable) BreakDownNames() []Name {
	var names []Name

	for _, s := range f.users {
		splits := strings.Split(s, ".")
		names = append(names, Name{splits[0], splits[1]})
	}
	return names
}
