package dom

import (
	"syscall/js"
)

// A GoNode is a node in the virtual DOM, representing either a DOM element or a text node.
// It may be used to create a (non-virtual) DOM node via BuildJsNode(), which returns
// a [JsNode] reference to the new browser DOM node.
// [GoElement] is a GoNode for a virtual DOM element.
// [GoTextNode] is a GoNode for a virtual DOM text node.
//
// The "virtual DOM" is a lightweight, in-memory representation of the
// browser's actual DOM (Document Object Model) elements.
type GoNode interface {
	BuildJsNode() JsNode       // builds and returns a browser DOM node wrapped in a JsNode
	BuildRawDomNode() js.Value // builds and returns a browser DOM node
	DebugName() string         // returns a string describing the node
	DebugJsValue() js.Value    // returns a JavaScript value describing the node
}

// JsNode is a strongly-typed js.Value referencing an actual browser DOM node,
// which is either a DOM element or text node.
// [JsElement] is a JsNode whose NodeJsValue() represents a browser DOM element.
// [JsTextNode] is a JsNode whose JsNodeValue() represents a browser DOM text node.
type JsNode interface {
	JsValue() js.Value
	isJsNode() // remove when JsNode gets anything distinctive
}

// jsNode is a JsNode whose nodeType is unspecified.
// When possible, you should use more accurate types like [JsElement] or [JsTextNode].
type jsNode struct {
	node js.Value
}

var _ JsNode = &jsNode{}

func (n *jsNode) isJsNode() {}
func (n *jsNode) JsValue() js.Value {
	return n.node
}

func NewNonspecificJsNode(node js.Value) JsNode {
	return &jsNode{node: node}
}

// BuildUntypedNodes builds DOM nodes for each element of the slice and returns
// them as []interface{}. Some DOM methods require []interface{}; for example,
// syscall/js.Value.Call().
func BuildUntypedNodes(nodes []GoNode) []interface{} {
	jsInterfaces := make([]interface{}, len(nodes))
	for i, node := range nodes {
		jsInterfaces[i] = node.BuildRawDomNode()
	}
	return jsInterfaces
}

func log(args ...interface{}) {
	js.Global().Get("console").Call("log", args...)
}
