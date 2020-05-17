package conthego

import "encoding/xml"

// https://stackoverflow.com/questions/30256729/how-to-traverse-through-xml-data-in-golang
type Command struct {
	node        *Node
	instruction string
	styles      string
}

func (c Command) success() {
	c.node.Attrs = append(c.node.Attrs, xml.Attr{xml.Name{"", "class"}, "success"})
}

func (c Command) failure() {
	c.node.Attrs = append(c.node.Attrs, xml.Attr{xml.Name{"", "class"}, "failure"})
}

func (c Command) echo(value string) {
	c.node.Content = value
}

func (c Command) getTextVal() string {
	return c.node.Content
}
