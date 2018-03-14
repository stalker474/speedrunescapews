package runescape

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

// RenderNode simple helper to convert a dom node to string
func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

// FindNode Recursively search for a node by type and id
// if id is "", only by type
func FindNode(n *html.Node, nodeType string, id string) *html.Node {
	if n.Type == html.ElementNode && n.Data == nodeType {
		if id != "" {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					return n
				}
			}
		} else {
			return n
		}
	}
	var foundNode *html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found := FindNode(c, nodeType, id)
		if found != nil {
			foundNode = found
		}
	}
	return foundNode
}

// FindNodeByClass Recursively search for a node by type and class
func FindNodeByClass(n *html.Node, nodeType string, class string) *html.Node {
	if n.Type == html.ElementNode && n.Data == nodeType {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == class {
				return n
			}
		}
	}
	var foundNode *html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found := FindNodeByClass(c, nodeType, class)
		if found != nil {
			foundNode = found
		}
	}
	return foundNode
}
