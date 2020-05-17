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
		if instr[0] == '?' && strings.HasSuffix(instr, ")") { // assert method call
			strValue := callMethod(f, instr[1:len(instr)], command.getTextVal())
			assertEquals(&command, command.getTextVal(), strValue)

		} else if instr[0] == '?' { // assert var
			strValue := f.getVar(instr[1:len(instr)])
			assertEquals(&command, command.getTextVal(), strValue)

		} else if instr[0] == '$' && strings.HasSuffix(instr, ")") { // echo method call
			strValue := callMethod(f, instr[1:len(instr)], command.getTextVal())
			command.echo(strValue)

		} else if instr[0] == '$' { // echo var
			command.echo(f.getVar(instr[1:len(instr)]))

		} else if strings.HasSuffix(instr, ")") { // var assignment, method call
			varName := instr[0:strings.Index(instr, "=")]
			methodCall := instr[strings.Index(instr, "=")+1 : len(instr)]
			strValue := callMethod(f, methodCall, command.getTextVal())
			f.putVar(varName, strValue)

		} else { // var assignment
			f.putVar(instr, command.node.Content)
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
