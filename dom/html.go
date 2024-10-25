package dom

import (
	"syscall/js"
)

// GoHtmlNode is a GoNode created from an HTML string.
type GoHtmlNode struct {
	html string
}

var _ GoNode = &GoHtmlNode{}

// NewGoHtmlNode creates a new GoHtmlNode from the provided HTML string.
func NewGoHtmlNode(html string) *GoHtmlNode {
	return &GoHtmlNode{html: html}
}

// DebugName returns a string representation of the GoHtmlNode for debugging purposes.
func (h *GoHtmlNode) DebugName() string {
	return "GoHtmlNode: `" + h.html + "`"
}

// DebugJsValue returns a js.Value representation of the GoHtmlNode for debugging purposes.
func (h *GoHtmlNode) DebugJsValue() js.Value {
	return js.ValueOf("HTML: `" + h.html + "`")
}

// BuildRawDomNode builds and returns a browser DOM node as a raw js.Value.
// If the HTML string renders multiple nodes, the content is wrapped in a DIV element.
// Use BuildRawDomNodes() to get the individual nodes.
func (h *GoHtmlNode) BuildRawDomNode() js.Value {
	div := js.Global().Get("document").Call("createElement", "div")
	div.Set("innerHTML", h.html)
	if div.Get("childNodes").Length() == 1 {
		return div.Get("firstChild")
	}
	return div
}

// BuildJsNode builds and returns a browser DOM node wrapped in a JsHtmlNode as a generic JsNode.
func (h *GoHtmlNode) BuildJsNode() JsNode {
	return h.BuildJsHtmlNode()
}

// BuildJsHtmlNode builds and returns a browser DOM node wrapped in a JsHtmlNode.
func (h *GoHtmlNode) BuildJsHtmlNode() *JsHtmlNode {
	return NewJsHtmlNode(h.BuildRawDomNode())
}

// JsHtmlNode is a strongly typed wrapper around a js.Value representing the DOM node created by an HTML string.
type JsHtmlNode struct {
	node js.Value
}

var _ JsNode = &JsHtmlNode{}

func (h JsHtmlNode) isJsNode() {}

func NewJsHtmlNode(htmlNode js.Value) *JsHtmlNode {
	if htmlNode.IsNull() || htmlNode.IsUndefined() {
		panic("NewJsHtmlNode: js.Value is null or undefined")
	}
	return &JsHtmlNode{node: htmlNode}
}

func (h JsHtmlNode) JsValue() js.Value {
	return h.node
}
