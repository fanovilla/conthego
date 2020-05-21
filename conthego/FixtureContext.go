package conthego

import (
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
	if strings.Contains(rawVar, "]") { // simple indexed
		open := strings.Index(rawVar, "[")
		close := strings.Index(rawVar, "]")
		rawVal := f.getVar(rawVar[0:open])
		if rawVal == nil {
			return "eval of " + rawVar[0:open] + " returned nil"
		}
		index, err := strconv.Atoi(rawVar[open+1 : close])
		if err != nil {
			panic(err)
		}
		// https://stackoverflow.com/questions/14025833/range-over-interface-which-stores-a-slice
		switch reflect.TypeOf(rawVal).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(rawVal)
			rawValAtIndex := formatAtom(s.Index(index))
			if strings.Contains(rawVar, ".") {
				keyString := rawVar[strings.Index(rawVar, ".")+1:]
				content, _ := dotnotation.Get(rawValAtIndex, keyString)
				return content
			}
			return rawValAtIndex
		}
		return "invalid index " + rawVar

	} else if strings.Contains(rawVar, ".") { // dot-notation
		rawVal := f.getVar(rawVar[0:strings.Index(rawVar, ".")])
		keyString := rawVar[strings.Index(rawVar, ".")+1:]
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
