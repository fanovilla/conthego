package conthego

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"io/ioutil"
	"os"
	//"os"
	"runtime"
	"strings"
)

func RunSpec() {
	baseName := getSpecBaseName()
	content := readFile(baseName + ".md")
	html := markdown.ToHTML(content, nil, nil)

	rootNode := unmarshal(html)
	commands := make([]Command, 0)
	collectCommands(rootNode, &commands)
	processCommands(&commands)

	bytes := marshal(rootNode)
	fmt.Println(string(bytes))

	writeFile(baseName, bytes)
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

//if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
//// path/to/whatever does not exist
//}

//n := new(big.Float)
//n, ok := n.SetString("44.3355")
//if !ok {
//	fmt.Println("SetString: error")
//	return
//}
