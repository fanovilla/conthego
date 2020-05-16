package conthego

type Fixture struct {
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
	for i := range *commands {
		command := (*commands)[i]
		instr := command.instruction
		if instr[0] == '?' {
			method := instr[1 : len(instr)-2]
			outs, err := InvokeMethod(Fixture{}, method)
			if err != nil {
				panic(err)
			}
			out := outs[0]
			atom := formatAtom(out)

			strValue, ok := atom.(string)
			if !ok {
				panic("not a string")
			}

			strExpected := command.node.Content
			if strExpected == strValue {
				command.success()
			} else {
				command.failure()
			}

		}
	}
}
