package examples

import (
	"github.com/fanovilla/conthego/conthego"
	"testing"
)

func TestPlotlyHack(t *testing.T) {
	conthego.RunSpec(t, &PlotlyHackFixture{})
}

type PlotlyHackFixture struct {
}

func (f *PlotlyHackFixture) HeadResources() []string {
	return []string{"concordion.css", "https://cdn.plot.ly/plotly-1.2.0.min.js"}
}

func (f *PlotlyHackFixture) EmbedRaw() string {
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
