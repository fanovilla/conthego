package conthego

import "strings"

type Fixture struct {
	vars map[string]string
}

func NewFixture() *Fixture {
	var f Fixture
	f.vars = make(map[string]string)
	return &f
}

func (f Fixture) putVar(name string, value string) {
	f.vars[name] = value
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

func processCommands(commands *[]Command) {
	f := NewFixture()
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
		} else {
			// variable assignment
			f.putVar(instr, command.node.Content)
		}
	}
}
