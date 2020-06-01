package examples

import (
	"encoding/xml"
	"github.com/fanovilla/conthego/conthego"
	"github.com/tidwall/gjson"
	"testing"
)

func TestEmbedRaw(t *testing.T) {
	conthego.RunSpec(t, &EmbedRawFixture{})
}

type EmbedRawFixture struct {
}

func (f *EmbedRawFixture) HeadResources() []string {
	return []string{"concordion.css", "https://cdn.plot.ly/plotly-1.2.0.min.js"}
}

func (f *EmbedRawFixture) EmbedTable() string {
	jsonString := `{
  "names" : [
    {
      "first": "mary",
      "last": "smith"
    },
    {
      "first": "lisa",
      "last": "brown"
    }
  ]
}`

	names := gjson.Get(jsonString, "names").Array()
	keys := []string{"first", "last"}
	t := newTable(keys)
	for i := 0; i < len(names); i++ {
		var elems []string
		for _, k := range keys {
			elems = append(elems, names[i].Get(k).String())
		}
		t.addRow(tr{Tds: elems})
	}

	bytes, _ := xml.Marshal(t)
	return string(bytes)
}

func (f *EmbedRawFixture) EmbedPlotly() string {
	return `
		<div id="tester" style="width:600px;height:250px;"></div>
		<script>
			TESTER = document.getElementById('tester');
			Plotly.newPlot( TESTER, [{
			x: [1, 2, 3, 4, 5],
			y: [1, 2, 4, 8, 16] }], {
			margin: { t: 0 } } );
		</script>
	`
}
