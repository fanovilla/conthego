package conthego

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"reflect"
	"strings"
)

func preProcess(node *html.Node) {
	if node.DataAtom == atom.Table {
		processTable(node)
	} else {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			preProcess(c)
		}
	}
}

func addHeaderResources(rootNode *html.Node, f *fixtureContext) {
	var headResources []string
	_, ok := reflect.TypeOf(f.localFixture).MethodByName("HeadResources")
	if ok {
		rawVal := callMethod(f, "HeadResources()", "")
		headResources = reflect.ValueOf(rawVal).Interface().([]string)
	}
	head := child(rootNode.FirstChild, atom.Head)
	removeStyle := false

	if headResources != nil {
		for _, r := range headResources {
			if strings.HasSuffix(r, ".css") {
				head.AppendChild(&html.Node{Type: html.ElementNode, DataAtom: atom.Link, Data: "link",
					Attr: []html.Attribute{attr("href", r), attr("rel", "stylesheet")},
				})
				removeStyle = true
			} else if strings.HasSuffix(r, ".js") {
				head.AppendChild(&html.Node{Type: html.ElementNode, DataAtom: atom.Script, Data: "script",
					Attr: []html.Attribute{attr("src", r)},
				})
			} else {
				fmt.Println("Unknown extension, not adding " + r)
			}
		}
	}

	if removeStyle {
		head.RemoveChild(child(head, atom.Style))
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
	thead := child(table, atom.Thead)
	if thead == nil {
		return
	}
	theadTr := child(thead, atom.Tr)
	ths := children(theadTr, atom.Th) // table>thead>tr>th
	for i := range ths {              // for each header col
		anchor := child(ths[i], atom.A) // find first command and check presence of iterator ":"
		if anchor != nil {
			instr := collectAttrs(anchor)["title"]
			if isCommand(anchor) && strings.Contains(instr, ":") {
				anchor.Parent.RemoveChild(anchor)
				anchor.Attr = []html.Attribute{attr("href", "-"), attr("title", strings.ReplaceAll(instr, ":", "="))}
				table.Parent.InsertBefore(anchor, table) //hoist up the iterator command; for convenience
				processTableStructs(table, instr[:strings.Index(instr, ":")])
			} else {
				processTablePerRow(table)
			}
		}
	}
}

func processTablePerRow(table *html.Node) {
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

func processTableStructs(table *html.Node, loopVar string) {
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
					text := ""
					if textNode != nil {
						text = textNode.Data
						textNode.Parent.RemoveChild(textNode)
					}
					tds[i].AppendChild(newAnchor(text, fmt.Sprintf("%s[%d]%s", pre, j, post)))
				}
			}
		}
	}

	rowMatcher := newAnchor(fmt.Sprintf("%d", len(trs)), "#"+loopVar)
	table.Parent.InsertBefore(rowMatcher, table.NextSibling)

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
