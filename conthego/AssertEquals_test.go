package conthego

import (
	"testing"
)

func Test(t *testing.T) {
	RunSpec()
}

func (f Fixture) Hello() string {
	return "World"
}
