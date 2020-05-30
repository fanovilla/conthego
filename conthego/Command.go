package conthego

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strconv"
)

// https://stackoverflow.com/questions/30256729/how-to-traverse-through-xml-data-in-golang
type Command struct {
	node        *html.Node
	instruction string
	styles      string
}

func (c Command) echo(value string, raw bool) {
	if raw {
		styleNode(c.node, "")
		n := html.Node{Type: html.RawNode, Data: value}
		c.node.Parent.InsertBefore(&n, c.node)
	} else {
		textNode := textNode(c.node)
		textNode.Data += value
		styleNode(c.node, "")
	}
}

func (c Command) getTextVal() string {
	return c.node.FirstChild.Data
}

func (c Command) success(f *fixtureContext) {
	if f.expectedToFail {
		styleNode(c.node, "failure")
		c.node.Data = fmt.Sprintf("%s (assert OK but expected to fail)", c.node.Data)
		f.t.Fail()
	} else {
		styleNode(c.node, "success")
	}
}

func (c Command) restyle() {
	styleNode(c.node, "")
}

func (c Command) failure(f *fixtureContext, actual string) {
	styleNode(c.node, "failure")
	c.node.FirstChild.Data = fmt.Sprintf("%s (actual=%s)", c.node.FirstChild.Data, actual)
	if !f.expectedToFail {
		f.t.Fail()
	}
}

func styleNode(node *html.Node, class string) {
	if class != "" {
		node.Attr = []html.Attribute{{Key: "class", Val: class}}
	} else {
		node.Attr = nil
	}
	node.DataAtom = atom.Span
	node.Data = "span"
}

func (c Command) assert(f *fixtureContext, val interface{}) {
	if val == nil {
		c.failure(f, "nil")
	} else if actual, ok := val.(bool); ok {
		if actual {
			c.success(f)
		} else {
			c.failure(f, "false")
		}
	} else if actual, ok := val.(float64); ok {
		actualStr := strconv.FormatFloat(actual, 'f', -1, 64)
		if c.node.FirstChild.Data == actualStr {
			c.success(f)
		} else {
			c.failure(f, actualStr)
		}
	} else if actual, ok := val.(string); ok {
		if c.node.FirstChild.Data == actual {
			c.success(f)
		} else {
			c.failure(f, actual)
		}
	} else {
		c.failure(f, "invalid assert not a bool or string value")
	}
}
