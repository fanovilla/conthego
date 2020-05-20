package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"sort"
	"strings"
	"testing"
)

func TestVerifyRows(t *testing.T) {
	conthego.RunSpec(t, &FixtureVerifyRows{})
}

type FixtureVerifyRows struct {
	users []string
}

func (f *FixtureVerifyRows) SetUpUser(user string) []string {
	f.users = append(f.users, user)
	return f.users
}

func (f *FixtureVerifyRows) SearchString(search string) []string {
	var results []string

	for _, s := range f.users {
		if strings.Contains(s, search) {
			results = append(results, s)
		}
	}
	sort.Strings(results)
	return results
}
