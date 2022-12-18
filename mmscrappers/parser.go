package mmscrappers

import (
	"bytes"
	"golang.org/x/net/html"
)

func getNodeTxt(n *html.Node) string {
	buf := bytes.Buffer{}
	collectText(n, &buf)
	return buf.String()
}

func collectText(n *html.Node, buf *bytes.Buffer) {

	if n == nil {
		return
	}

	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

func getAttr(key string, n html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}

	return ""
}
