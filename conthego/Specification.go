package conthego

import (
	"encoding/xml"
	"fmt"
	"github.com/gomarkdown/markdown"
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

	rootNode := unmarshal(html)
	runCommands(rootNode, f)

	bytes := marshal(rootNode)
	fmt.Println(string(bytes))
	writeFile(baseName, bytes)
}

func runCommands(rootNode *Node, f *fixtureContext) {
	commands := make([]Command, 0)
	preProcess(rootNode)

	bytes := marshal(rootNode)
	fmt.Println("after pre-processing:" + string(bytes))

	collectCommands(rootNode, &commands)
	reportLines := processCommands(f, &commands)

	reportNode := Node{xml.Name{Local: "div"}, []xml.Attr{}, "", []Node{}}
	for _, s := range reportLines {
		reportNode.Nodes = append(reportNode.Nodes, Node{xml.Name{Local: "p"}, []xml.Attr{}, s, []Node{}})
	}
	rootNode.Nodes[1].Nodes = append(rootNode.Nodes[1].Nodes, reportNode)
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
