package conthego

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

// https://stackoverflow.com/questions/30256729/how-to-traverse-through-xml-data-in-golang
type Command struct {
	node        *Node
	instruction string
	styles      string
}

func (c Command) echo(value string) {
	c.node.Content = value
	styleNode(c.node, "")
}

func (c Command) getTextVal() string {
	return c.node.Content
}

func (c Command) success() {
	styleNode(c.node, "success")
}

func (c Command) restyle() {
	styleNode(c.node, "")
}

func (c Command) failure(f *fixtureContext, actual string) {
	styleNode(c.node, "failure")
	c.node.Content = fmt.Sprintf("%s (actual=%s)", c.node.Content, actual)
	if !f.expectedToFail {
		f.t.Fail()
	}
}

func styleNode(node *Node, class string) {
	attrs := []xml.Attr{}
	for i := range node.Attrs {
		switch name := node.Attrs[i].Name.Local; name {
		case "href", "title":
			// NOOP
		default:
			attrs = append(attrs, node.Attrs[i])
		}
	}
	if class != "" {
		attrs = append(attrs, xml.Attr{xml.Name{"", "class"}, class})
	}
	node.Attrs = attrs
	node.XMLName.Local = "span"
}

func (c Command) assert(f *fixtureContext, val interface{}) {
	if val == nil {
		c.failure(f, "nil")
	} else if actual, ok := val.(bool); ok {
		if actual {
			c.success()
		} else {
			c.failure(f, "false")
		}
	} else if actual, ok := val.(float64); ok {
		actualStr := strconv.FormatFloat(actual, 'f', -1, 64)
		if c.node.Content == actualStr {
			c.success()
		} else {
			c.failure(f, actualStr)
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
