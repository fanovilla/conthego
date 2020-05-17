package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestReturnMap(t *testing.T) {
	conthego.RunSpec(conthego.NewFixture(t, &FixtureReturnMap{}))
}

type FixtureReturnMap struct {
}

type Vertex struct {
	Lat, Long float64
}

func (f FixtureReturnMap) BuildLocationMap() map[string]Vertex {
	return map[string]Vertex{
		"BellLabs": Vertex{
			40.68433, -74.39967,
		},
		"Google": Vertex{
			37.42202, -122.08408,
		},
	}
}
