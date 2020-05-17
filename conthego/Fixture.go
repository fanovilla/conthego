package conthego

import (
	"fmt"
	"github.com/joeycumines/go-dotnotation/dotnotation"
	"strings"
	"testing"
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
	if strings.Contains(rawVar, ".") { // dot-notation
		rawVal := f.getVar(rawVar[0:strings.Index(rawVar, ".")])
		keyString := rawVar[strings.Index(rawVar, ".")+1 : len(rawVar)]
		content, _ := dotnotation.Get(rawVal, keyString)
		return content
	} else {
		return f.getVar(rawVar)
	}
}

func collectCommands(node *Node, commands *[]Command) {
	m := make(map[string]string)
	for _, a := range node.Attrs {
		m[a.Name.Local] = a.Value
	}

	if m["href"] == "-" {
		*commands = append(*commands, Command{node, m["title"], m["styles"]})
	}

	for i := range node.Nodes {
		collectCommands(&(node.Nodes[i]), commands)
	}
}

func processCommands(f *fixtureContext, commands *[]Command) {
	for i := range *commands {
		command := (*commands)[i]
		instr := command.instruction
		if instr[0] == '!' {
			// directive
			if "ExpectedToFail" == instr[1:len(instr)] {
				f.expectedToFail = true
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

		} else if strings.HasSuffix(instr, ")") {
			// var assignment, method call
			varName := instr[0:strings.Index(instr, "=")]
			methodCall := instr[strings.Index(instr, "=")+1 : len(instr)]
			strValue := callMethod(f, methodCall, command.getTextVal())
			f.putVar(varName, strValue)

		} else {
			// var assignment
			f.putVar(instr, command.getTextVal())
		}
	}
}
