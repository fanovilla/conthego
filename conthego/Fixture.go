package conthego

import (
	"github.com/joeycumines/go-dotnotation/dotnotation"
	"strings"
)

type fixtureContext struct {
	vars         map[string]interface{}
	localFixture interface{}
}

func NewFixture(internalFixture interface{}) *fixtureContext {
	var f fixtureContext
	f.vars = make(map[string]interface{})
	f.localFixture = internalFixture
	return &f
}

func (f fixtureContext) putVar(name string, value interface{}) {
	f.vars[name] = value
}

func (f fixtureContext) getVar(name string) interface{} {
	return f.vars[name]
}

func (f fixtureContext) evalVar(rawVar string) string {
	var strValue string
	if strings.Contains(rawVar, ".") { // dot-notation
		rawVal := f.getVar(rawVar[0:strings.Index(rawVar, ".")])
		keyString := rawVar[strings.Index(rawVar, ".")+1 : len(rawVar)]
		content, _ := dotnotation.Get(rawVal, keyString)
		if content == nil {
			strValue = "nil"
		} else {
			strValue = content.(string)
		}
	} else {
		strValue = f.getVar(rawVar).(string)
	}
	return strValue
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
		if instr[0] == '?' && strings.HasSuffix(instr, ")") {
			// assert method call
			strValue := callMethod(f, instr[1:len(instr)], command.getTextVal())
			assertEquals(&command, command.getTextVal(), strValue.(string))

		} else if instr[0] == '?' {
			// assert var
			strValue := f.evalVar(instr[1:len(instr)])
			assertEquals(&command, command.getTextVal(), strValue)

		} else if instr[0] == '$' && strings.HasSuffix(instr, ")") {
			// echo method call
			strValue := callMethod(f, instr[1:len(instr)], command.getTextVal())
			command.echo(strValue.(string))

		} else if instr[0] == '$' {
			// echo var
			strValue := f.evalVar(instr[1:len(instr)])
			command.echo(strValue)

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

func assertEquals(command *Command, expected string, actual string) {
	if expected == actual {
		command.success()
	} else {
		command.failure()
	}
}
