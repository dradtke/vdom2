package vdom

import (
	"bytes"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"honnef.co/go/js/dom"
)

// Apply updates
func Apply(el dom.Element, src []byte) {
	doApply(el, parse(src))
}

func doApply(el dom.Element, want *html.Node) {
	haveChild, wantChild := el.FirstChild(), want.FirstChild
	for haveChild != nil || wantChild != nil {
		switch {
		case haveChild == nil:
			// append
			el.AppendChild(createNode(wantChild))
			wantChild = wantChild.NextSibling

		case wantChild == nil:
			// delete
			el.RemoveChild(haveChild)
			haveChild = haveChild.NextSibling()

		case nodesDiffer(haveChild, wantChild):
			// modify
			// TODO: this can likely be optimized based on what needs to change
			newChild := createNode(wantChild)
			el.InsertBefore(newChild, haveChild)
			el.RemoveChild(haveChild)
			haveChild, wantChild = newChild.NextSibling(), wantChild.NextSibling

		default:
			// these nodes are the same type, so diff recursively
			// if they're elements, or modify the content if
			// they're text nodes
			switch t := haveChild.(type) {
			case dom.Element:
				doApply(t, wantChild)

			case *dom.Text:
				if t.TextContent() != wantChild.Data {
					t.SetTextContent(wantChild.Data)
				}
			}
			haveChild, wantChild = haveChild.NextSibling(), wantChild.NextSibling
		}
	}
}

func parse(src []byte) *html.Node {
	nodes, err := html.ParseFragment(bytes.NewReader(src), &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
	})
	if err != nil {
		panic("failed to parse HTML tree: " + err.Error())
	}
	return nodes[0]
}

func createNode(node *html.Node) dom.Node {
	// TODO: flesh this out more
	switch node.Type {
	case html.ElementNode:
		return dom.GetWindow().Document().CreateElement(node.Data)

	case html.TextNode:
		return dom.GetWindow().Document().CreateTextNode(node.Data)

	default:
		panic("can't create node")
	}
}

var domToHtmlNodeType = map[int]html.NodeType{
	1: html.ElementNode,
	3: html.TextNode,
	// 7: processing_instruction_node,
	8:  html.CommentNode,
	9:  html.DocumentNode,
	10: html.DoctypeNode,
	// 11: document_fragment,
}

func nodesDiffer(x dom.Node, y *html.Node) bool {
	return domToHtmlNodeType[x.NodeType()] != y.Type
}
