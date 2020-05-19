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
}

var users []string

func (f FixtureVerifyRows) SetUpUser(user string) []string {
	users = append(users, user)
	return users
}

func (f FixtureVerifyRows) SearchString(search string) []string {
	var results []string

	for _, s := range users {
		if strings.Contains(s, search) {
			results = append(results, s)
		}
	}
	sort.Strings(results)
	return results
}
