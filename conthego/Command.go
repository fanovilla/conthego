package conthego

import (
	"encoding/xml"
	"fmt"
)

// https://stackoverflow.com/questions/30256729/how-to-traverse-through-xml-data-in-golang
type Command struct {
	node        *Node
	instruction string
	styles      string
}

func (c Command) echo(value string) {
	c.node.Content = value
}

func (c Command) getTextVal() string {
	return c.node.Content
}

func (c Command) assert(f *fixtureContext, actual string) {
	if c.node.Content == actual {
		c.node.Attrs = append(c.node.Attrs, xml.Attr{xml.Name{"", "class"}, "success"})
	} else {
		c.node.Attrs = append(c.node.Attrs, xml.Attr{xml.Name{"", "class"}, "failure"})
		c.node.Content = fmt.Sprintf("%s (actual=%s)", c.node.Content, actual)
		if !f.expectedToFail {
			f.t.Fail()
		}
	}
}
