package conthego

import (
	"encoding/xml"
	"fmt"
	"strings"
)

var verifyRows = false

func preProcess(node *Node) {
	singleUseVerifyRows := verifyRows
	verifyRows = false
	if node.elem() == "table" {
		if singleUseVerifyRows {
			processTableStructs(node)
		} else {
			processTable(node)
		}
	} else if node.elem() == "ul" && singleUseVerifyRows {
		processListStrings(node)
	} else if node.elem() == "a" { // check for directives
		m := collectAttrs(node)
		if m["title"] == "!VerifyRows" {
			verifyRows = true
			node.Attrs = nil // consume directive
			node.XMLName.Local = "span"
		}
	} else {
		for i := range node.Nodes {
			preProcess(&(node.Nodes[i]))
		}
	}
}

func processTable(table *Node) {
	ths := table.Nodes[0].Nodes[0].Nodes // table>thead>tr>th
	trs := table.Nodes[1].Nodes          // table>tbody>tr
	for i := range ths {                 // for each header col
		if len(ths[i].Nodes) > 0 { // assume has anchor child
			anchor := &(ths[i].Nodes[0])
			m := collectAttrs(anchor)
			if m["href"] == "-" { // if header col has command copy it to row col
				for j := range trs { // for each row
					trs[j].Nodes[i].Nodes = []Node{*anchor}
					trs[j].Nodes[i].Nodes[0].Content = trs[j].Nodes[i].Content
					trs[j].Nodes[i].Content = ""
				}
				removeAnchor(anchor)
			}
		}
	}
}

func processTableStructs(table *Node) {
	ths := table.Nodes[0].Nodes[0].Nodes // table>thead>tr>th
	trs := table.Nodes[1].Nodes          // table>tbody>tr
	for i := range ths {                 // for each header col
		if len(ths[i].Nodes) > 0 { // assume has anchor child
			anchor := &(ths[i].Nodes[0])
			m := collectAttrs(anchor)
			if m["href"] == "-" { // if header col has command copy it to row col
				rawTitle := m["title"]
				dotPos := strings.Index(rawTitle, ".")
				if dotPos >= 0 {
					pre := rawTitle[0:dotPos]
					post := rawTitle[dotPos:len(rawTitle)]
					for j := range trs { // for each row
						a := newAnchor(trs[j].Nodes[i].Content, fmt.Sprintf("%s[%d]%s", pre, j, post))
						trs[j].Nodes[i].Nodes = []Node{*a}
						trs[j].Nodes[i].Content = ""
					}
				}
				removeAnchor(anchor)
			}
		}
	}
}

func newAnchor(content string, title string) *Node {
	hrefAttr := attr("href", "-")
	titleAttr := attr("title", title)
	return &Node{
		XMLName: xml.Name{
			Local: "a",
		},
		Attrs:   []xml.Attr{hrefAttr, titleAttr},
		Content: content,
		Nodes:   nil,
	}
}

func removeAnchor(node *Node) { // remove command that was consumed at pre-process step
	if node.XMLName.Local == "a" {
		node.XMLName.Local = "span"
		node.Attrs = nil
	}
}

func attr(name string, val string) xml.Attr {
	return xml.Attr{
		Name: xml.Name{
			Local: name,
		},
		Value: val,
	}
}

func processListStrings(ul *Node) {
	lis := ul.Nodes            // ul>li
	if len(lis[0].Nodes) > 0 { // assume command exists
		anchor := lis[0].Nodes[0]
		lis[0].Content = anchor.Content
		m := collectAttrs(&anchor)
		for i := range lis { // for each row
			a := newAnchor(lis[i].Content, fmt.Sprintf("%s[%d]", m["title"], i))
			lis[i].Nodes = []Node{*a}
			lis[i].Content = ""
		}
	}
}
