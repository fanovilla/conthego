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

func (c Command) success() {
	c.node.Attrs = append(c.node.Attrs, xml.Attr{xml.Name{"", "class"}, "success"})
}

func (c Command) failure(f *fixtureContext, actual string) {
	c.node.Attrs = append(c.node.Attrs, xml.Attr{xml.Name{"", "class"}, "failure"})
	c.node.Content = fmt.Sprintf("%s (actual=%s)", c.node.Content, actual)
	if !f.expectedToFail {
		f.t.Fail()
	}
}

func (c Command) assert(f *fixtureContext, val interface{}) {
	if actual, ok := val.(bool); ok {
		if actual {
			c.success()
		} else {
			c.failure(f, "false")
		}
	} else if actual, ok := val.(string); ok {
		if c.node.Content == actual {
			c.success()
		} else {
			c.failure(f, actual)
		}
	} else {
		c.failure(f, "invalid assert not a bool or string value")
	}
}
