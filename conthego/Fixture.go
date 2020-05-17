package conthego

import "strings"

type fixtureContext struct {
	vars         map[string]string
	localFixture interface{}
}

func NewFixture(internalFixture interface{}) *fixtureContext {
	var f fixtureContext
	f.vars = make(map[string]string)
	f.localFixture = internalFixture
	return &f
}

func (f fixtureContext) putVar(name string, value string) {
	f.vars[name] = value
}

func (f fixtureContext) getVar(name string) string {
	return f.vars[name]
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
		if instr[0] == '?' && strings.HasSuffix(instr, ")") { // method call
			strValue := callMethod(f, instr[1:len(instr)])
			strExpected := command.node.Content
			if strExpected == strValue {
				command.success()
			} else {
				command.failure()
			}
		} else if instr[0] == '$' && strings.HasSuffix(instr, ")") { // echo method call
			strValue := callMethod(f, instr[1:len(instr)])
			command.echo(strValue)
		} else if instr[0] == '$' { // echo var
			command.echo(f.getVar(instr[1:len(instr)]))
		} else {
			// variable assignment
			f.putVar(instr, command.node.Content)
		}
	}
}
