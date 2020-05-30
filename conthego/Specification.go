package conthego

import (
	"bytes"
	"fmt"
	"github.com/gomarkdown/markdown"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io/ioutil"
	"os"
	"runtime/debug"
	"testing"

	//"os"
	"runtime"
	"strings"
)

func RunSpec(t *testing.T, internalFixture interface{}) {
	// https://blog.golang.org/defer-panic-and-recover
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			t.Fatal(fmt.Sprint(r))
		}
	}()

	f := newFixture(t, internalFixture)
	baseName := getSpecBaseName()
	content := readFile(baseName + ".md")
	html := markdown.ToHTML(content, nil, nil)

	rootNode := unmarshalSpec(html)

	defer func() {
		if r := recover(); r != nil {
			appendStack(rootNode, baseName, fmt.Sprint(r))
			debug.PrintStack()
			t.Fatal(fmt.Sprint(r))
		}
	}()

	addHeaderResources(rootNode, f)
	runCommands(rootNode, f)

	bytes := marshalSpec(rootNode)
	fmt.Println(string(bytes))
	writeFile(baseName, bytes)
}

func appendStack(rootNode *html.Node, baseName string, error string) {
	stackNode := html.Node{Type: html.ElementNode, DataAtom: atom.Div, Data: "pre", Attr: []html.Attribute{attr("class", "failure")}}
	stackNode.AppendChild(&html.Node{Type: html.TextNode, Data: string(debug.Stack())})
	stackNode.AppendChild(&html.Node{Type: html.TextNode, Data: error})
	child(rootNode.FirstChild, atom.Body).AppendChild(&stackNode)

	bytes := marshalSpec(rootNode)
	fmt.Println(string(bytes))
	writeFile(baseName, bytes)
}

func runCommands(rootNode *html.Node, f *fixtureContext) {
	preProcess(rootNode)

	bytes := marshalSpec(rootNode)
	fmt.Println("after pre-processing:" + string(bytes))

	commands := make([]Command, 0)
	collectCommands(rootNode, &commands)
	reportLines := processCommands(f, &commands)

	reportNode := html.Node{Type: html.ElementNode, DataAtom: atom.Div, Data: "div", Attr: []html.Attribute{attr("class", "footer")}}
	for _, s := range reportLines {
		reportNode.AppendChild(newParagraph(s))
	}
	child(rootNode.FirstChild, atom.Body).AppendChild(&reportNode)
}

func newParagraph(text string) *html.Node {
	n := html.Node{Type: html.ElementNode, DataAtom: atom.P, Data: "p"}
	n.AppendChild(&html.Node{Type: html.TextNode, Data: text})
	return &n
}

func getSpecBaseName() string {
	fixtureFileName := ""
	for i := 0; fixtureFileName == "" && i < 3; i++ {
		_, fileName, _, ok := runtime.Caller(i)
		if !ok {
			panic("error finding fixture name with suffix '_test.go'")
		} else if strings.HasSuffix(fileName, "_test.go") {
			fixtureFileName = fileName
		}
	}
	return strings.TrimSuffix(fixtureFileName, "_test.go")
}

func readFile(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic("specification not found:" + filePath)
	}
	return content
}

func writeFile(filePath string, content []byte) {
	wd, _ := os.Getwd()
	name := wd + "/_" + filePath[strings.LastIndex(filePath, "/")+1:len(filePath)] + ".html"
	fmt.Println("Writing output file" + name)
	err := ioutil.WriteFile(name, content, 0666)
	if err != nil {
		panic("error writing output:" + filePath)
	}
}

func marshalSpec(rootNode *html.Node) []byte {
	var buf bytes.Buffer
	html.Render(&buf, rootNode)
	return buf.Bytes()
}

func unmarshalSpec(xhtml []byte) *html.Node {
	style := `
<style>
* {  font-family: Arial;}
pre { padding: 6px 28px 6px 28px;    background-color: #E8EEF7;}
pre, pre *, code, code *, kbd {  font-family: Courier New, Courier;  font-size: 10pt;}
table { border-collapse: collapse; empty-cells: show; margin: 8px 0px 8px 0px; }
th, td {  border: 1px solid black;  padding: 3px;}
th {  background-color: #C3D9FF;}
p, td, th, li, .breadcrumbs {  font-size: 10pt;}
.success, .success * {  background-color: #afa !important;}
.failure, .failure * {  background-color: #ffb0b0;  padding: 1px;}
.footer {  text-align: right;  margin-top: 40px;  font-size: 8pt;  width: 100%;  color: #999; }
</style>
`
	full := fmt.Sprintf("<head>%s</head><body>%s</body>", style, xhtml)
	rootNode, err := html.Parse(strings.NewReader(full))
	if err != nil {
		panic(err)
	}
	return rootNode
}
