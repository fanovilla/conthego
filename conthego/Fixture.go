package conthego

import (
	"encoding/xml"
	"fmt"
	"github.com/joeycumines/go-dotnotation/dotnotation"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

type fixtureContext struct {
	vars           map[string]interface{}
	localFixture   interface{}
	expectedToFail bool
	t              *testing.T
}

func newFixture(t *testing.T, internalFixture interface{}) *fixtureContext {
	var f fixtureContext
	f.vars = make(map[string]interface{})
	f.localFixture = internalFixture
	f.t = t
	return &f
}

func (f fixtureContext) putVar(name string, value interface{}) {
	f.vars[name] = value
}

func (f fixtureContext) getVar(name string) interface{} {
	return f.vars[name]
}

func (f fixtureContext) evalVar(rawVar string) interface{} {
	if strings.HasSuffix(rawVar, "]") { // simple indexed
		open := strings.Index(rawVar, "[")
		close := strings.Index(rawVar, "]")
		rawVal := f.getVar(rawVar[0:open])
		index, err := strconv.Atoi(rawVar[open+1 : close])
		if err != nil {
			panic(err)
		}
		// https://stackoverflow.com/questions/14025833/range-over-interface-which-stores-a-slice
		switch reflect.TypeOf(rawVal).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(rawVal)
			return formatAtom(s.Index(index))
		}
		return "invalid index " + rawVar

	} else if strings.Contains(rawVar, ".") { // dot-notation
		rawVal := f.getVar(rawVar[0:strings.Index(rawVar, ".")])
		keyString := rawVar[strings.Index(rawVar, ".")+1 : len(rawVar)]
		content, _ := dotnotation.Get(rawVal, keyString)
		return content
	} else {
		return f.getVar(rawVar)
	}
}

func collectCommands(node *Node, commands *[]Command) {
	m := collectAttrs(node)
	if m["href"] == "-" {
		*commands = append(*commands, Command{node, strings.TrimSpace(m["title"]), m["styles"]})
	}

	for i := range node.Nodes {
		collectCommands(&(node.Nodes[i]), commands)
	}
}

func collectAttrs(anchor *Node) map[string]string {
	m := make(map[string]string)
	for _, a := range anchor.Attrs {
		m[a.Name.Local] = a.Value
	}
	return m
}

func processTable(table *Node) {
	ths := table.Nodes[0].Nodes[0].Nodes // table>thead>tr>th
	trs := table.Nodes[1].Nodes          // table>tbody>tr
	for i := range trs {                 // for each row
		for j := range ths { // for each header col
			if len(ths[j].Nodes) > 0 { // assume has anchor child
				anchor := ths[j].Nodes[0]
				m := collectAttrs(&anchor)
				if m["href"] == "-" { // if header col has command copy it to row col
					trs[i].Nodes[j].Nodes = []Node{anchor}
					trs[i].Nodes[j].Nodes[0].Content = trs[i].Nodes[j].Content
					trs[i].Nodes[j].Content = ""
				}
			}
		}
	}
	for j := range ths { // for each header col
		if len(ths[j].Nodes) > 0 { // assume has anchor child
			ths[j].Nodes[0].Attrs = []xml.Attr{}
		}
	}
}

func attr(name string, val string) xml.Attr {
	return xml.Attr{
		Name: xml.Name{
			Local: name,
		},
		Value: val,
	}
}

func processListStrings(ul *Node) {
	lis := ul.Nodes            // ul>li
	if len(lis[0].Nodes) > 0 { // assume command exists
		anchor := lis[0].Nodes[0]
		lis[0].Content = anchor.Content
		m := collectAttrs(&anchor)
		for i := range lis { // for each row
			lis[i].Nodes = []Node{{
				XMLName: anchor.XMLName,
				Attrs:   []xml.Attr{attr("href", "-"), attr("title", fmt.Sprintf("%s[%d]", m["title"], i))},
				Content: lis[i].Content,
				Nodes:   nil,
			}}
			lis[i].Content = ""
		}
	}
}

var verifyRows = false

func preProcess(node *Node) {
	verifyRowsCalled := verifyRows
	verifyRows = false
	if node.elem() == "table" {
		if verifyRowsCalled {
			//TODO
		} else {
			processTable(node)
		}
	} else if node.elem() == "ul" && verifyRowsCalled {
		processListStrings(node)
	} else if node.elem() == "a" { // check for directives
		m := collectAttrs(node)
		if m["title"] == "!VerifyRows" {
			verifyRows = true
			node.Attrs = nil // consume directive
			node.XMLName.Local = "span"
		}
	} else {
		for i := range node.Nodes {
			preProcess(&(node.Nodes[i]))
		}
	}
}

func processCommands(f *fixtureContext, commands *[]Command) []string {
	var reportLines []string
	for i := range *commands {
		command := (*commands)[i]
		instr := command.instruction
		if instr[0] == '!' {
			// directive
			if "ExpectedToFail" == instr[1:len(instr)] {
				f.expectedToFail = true
				reportLines = append(reportLines, "This specification is ExpectedToFail")
			}

		} else if instr[0] == '?' && strings.HasSuffix(instr, ")") {
			// assert method call
			genericVal := callMethod(f, instr[1:len(instr)], command.getTextVal())
			command.assert(f, genericVal)

		} else if instr[0] == '?' {
			// assert var
			genericVal := f.evalVar(instr[1:len(instr)])
			command.assert(f, genericVal)

		} else if instr[0] == '$' && strings.HasSuffix(instr, ")") {
			// echo method call
			genericVal := callMethod(f, instr[1:len(instr)], command.getTextVal())
			command.echo(fmt.Sprint(genericVal))

		} else if instr[0] == '$' {
			// echo var
			genericVal := f.evalVar(instr[1:len(instr)])
			command.echo(fmt.Sprint(genericVal))

		} else if strings.HasSuffix(instr, ")") && strings.Contains(instr, "=") {
			// var assignment, method call
			varName := strings.TrimSpace(instr[0:strings.Index(instr, "=")])
			methodCall := strings.TrimSpace(instr[strings.Index(instr, "=")+1 : len(instr)])
			strValue := callMethod(f, methodCall, command.getTextVal())
			f.putVar(varName, strValue)
			command.restyle()

		} else if strings.HasSuffix(instr, ")") {
			// method call, no var assignment (fixture side-effect)
			methodCall := strings.TrimSpace(instr[strings.Index(instr, "=")+1 : len(instr)])
			callMethod(f, methodCall, command.getTextVal())
			command.restyle()

		} else {
			// var assignment
			f.putVar(instr, command.getTextVal())
			command.restyle()
		}
	}

	reportLines = append(reportLines, "Generated "+time.Now().Format(time.RFC1123Z))
	return reportLines
}
