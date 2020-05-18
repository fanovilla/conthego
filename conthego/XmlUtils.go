package conthego

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",attr"`
	Content string     `xml:",chardata"`
	Nodes   []Node     `xml:",any"`
}

func (n Node) elem() string {
	return n.XMLName.Local
}

// https://stackoverflow.com/questions/52669545/round-trip-xml-through-unmarshal-and-marshalindent
var replacer = strings.NewReplacer("&#xA;", "", "&#x9;", "", "\n", "", "\t", "")

func recursiveReplace(n *Node) {
	n.Content = replacer.Replace(n.Content)
	for i := range n.Nodes {
		recursiveReplace(&n.Nodes[i])
	}
}

func marshal(rootNode *Node) []byte {
	output, err := xml.Marshal(&rootNode)
	if err != nil {
		panic(err)
	}
	return output
}

func unmarshal(html []byte) *Node {
	rootNode := Node{}
	style := "<style>.success {background-color: #afa;} .failure {background-color: #ffb0b0;}</style>"
	err := xml.Unmarshal([]byte(fmt.Sprintf("<html><header>%s</header><body>%s</body></html>", style, html)), &rootNode)
	if err != nil {
		panic(err)
	}
	recursiveReplace(&rootNode)
	return &rootNode
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node Node
	return d.DecodeElement((*node)(n), &start)
}

func walk(nodes []Node, f func(Node) bool) {
	for _, n := range nodes {
		if f(n) {
			walk(n.Nodes, f)
		}
	}
}
