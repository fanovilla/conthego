package conthego

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
)

var verifyRows = false

func preProcess(node *html.Node) {
	singleUseVerifyRows := verifyRows
	verifyRows = false
	if node.DataAtom == atom.Table {
		if singleUseVerifyRows {
			processTableStructs(node)
		} else {
			processTable(node)
		}
	} else if node.DataAtom == atom.Ul && singleUseVerifyRows {
		processListStrings(node)
	} else if node.DataAtom == atom.A { // check for directives
		m := collectAttrs(node)
		if m["title"] == "!VerifyRows" {
			verifyRows = true
			node.Attr = nil // consume directive
			node.DataAtom = atom.Span
			node.Data = "span"
		}
	} else {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			preProcess(c)
		}
	}
}

func children(node *html.Node, tag atom.Atom) []*html.Node {
	var nodes []*html.Node
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		if n.DataAtom == tag {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func child(node *html.Node, tag atom.Atom) *html.Node {
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		if n.DataAtom == tag {
			return n
		}
	}
	return nil
}

func textNode(node *html.Node) *html.Node {
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.TextNode {
			return n
		}
	}
	return nil
}

func isCommand(n *html.Node) bool {
	return n.Type == html.ElementNode && collectAttrs(n)["href"] == "-"
}

func processTable(table *html.Node) {
	theadTr := child(child(table, atom.Thead), atom.Tr)
	ths := children(theadTr, atom.Th)                  // table>thead>tr>th
	trs := children(child(table, atom.Tbody), atom.Tr) // table>tbody>tr
	for j := range trs {                               // for each row
		tds := children(trs[j], atom.Td) // table>tbody>tr
		for i := range ths {             // for each header col
			anchor := child(ths[i], atom.A)
			if isCommand(anchor) {
				textNode := textNode(tds[i])
				tds[i].AppendChild(newAnchorFromAttrs(textNode.Data, anchor.Attr))
				textNode.Parent.RemoveChild(textNode)
			}
		}
	}
	for i := range ths { // for each header col
		anchor := child(ths[i], atom.A)
		if isCommand(anchor) {
			removeAnchor(anchor)
		}
	}
}

func processTableStructs(table *html.Node) {
	theadTr := child(child(table, atom.Thead), atom.Tr)
	ths := children(theadTr, atom.Th)                  // table>thead>tr>th
	trs := children(child(table, atom.Tbody), atom.Tr) // table>tbody>tr
	for j := range trs {                               // for each row
		tds := children(trs[j], atom.Td) // table>tbody>tr
		for i := range ths {             // for each header col
			anchor := child(ths[i], atom.A)
			if isCommand(anchor) {
				m := collectAttrs(anchor)
				rawTitle := m["title"]
				dotPos := strings.Index(rawTitle, ".")
				if dotPos >= 0 {
					pre := rawTitle[0:dotPos]
					post := rawTitle[dotPos:len(rawTitle)]
					textNode := textNode(tds[i])
					tds[i].AppendChild(newAnchor(textNode.Data, fmt.Sprintf("%s[%d]%s", pre, j, post)))
					textNode.Parent.RemoveChild(textNode)
				}
			}
		}
	}
	for i := range ths { // for each header col
		anchor := child(ths[i], atom.A)
		if isCommand(anchor) {
			removeAnchor(anchor)
		}
	}
}

func newAnchor(content string, title string) *html.Node {
	return newAnchorFromAttrs(content, []html.Attribute{attr("href", "-"), attr("title", title)})
}

func newAnchorFromAttrs(content string, attrs []html.Attribute) *html.Node {
	newAnchor := html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.A,
		Data:     "a",
		Attr:     attrs}
	newAnchor.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: content,
	})
	return &newAnchor
}

func removeAnchor(node *html.Node) { // remove command that was consumed at pre-process step
	if node.DataAtom == atom.A {
		node.DataAtom = atom.Span
		node.Data = "span"
		node.Attr = nil
	}
}

func attr(key string, val string) html.Attribute {
	return html.Attribute{
		Key: key,
		Val: val,
	}
}

func processListStrings(ul *html.Node) {
	//lis := ul.Nodes            // ul>li
	//if len(lis[0].Nodes) > 0 { // assume command exists
	//	anchor := lis[0].Nodes[0]
	//	lis[0].Content = anchor.Content
	//	m := collectAttrs(&anchor)
	//	for i := range lis { // for each row
	//		a := newAnchor(lis[i].Content, fmt.Sprintf("%s[%d]", m["title"], i))
	//		lis[i].Nodes = []Node{*a}
	//		lis[i].Content = ""
	//	}
	//}
}
