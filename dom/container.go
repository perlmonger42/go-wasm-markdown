package dom

import (
	"syscall/js"
)

// A JsNodeContainer is a JsNode that can contain other JsNodes.
type JsNodeContainer interface {
	JsNode
	AppendChild(child GoNode)
	AppendChildren(children ...GoNode)
	QuerySelector(selector string) (node JsNode, ok bool)
	QuerySelectorAll(selector string) *JsNodeList
}

// AppendChild appends a child node to a parent node.
func AppendChild(parent JsNodeContainer, child JsNode) {
	parent.JsValue().Call("appendChild", child.JsValue())
}

func AppendChildren(parent JsNodeContainer, children ...GoNode) {
	AppendChild(parent, NewGoFragment("temp", children...).BuildJsFragment())
}

// QuerySelector returns the first Element inside parent that matches the given CSS selector.
func QuerySelector(parent JsNodeContainer, selector string) (node JsNode, ok bool) {
	jsValue := parent.JsValue().Call("querySelector", selector)
	if jsValue.IsNull() || jsValue.IsUndefined() {
		return nil, false
	}
	return NewJsElement(jsValue), true
}

// QuerySelectorAll returns all Elements inside parent that match the given CSS selector.
// If no elements match, the list is empty.
func QuerySelectorAll(parent JsNodeContainer, selector string) *JsNodeList {
	jsValue := parent.JsValue().Call("querySelectorAll", selector)
	return NewJsNodeList(jsValue)
}

// JsNodeList represents a collection of nodes.
// It wraps a js.Value implementing the DOM NodeList interface.
type JsNodeList struct {
	jsValue js.Value
}

func NewJsNodeList(l js.Value) *JsNodeList {
	return &JsNodeList{jsValue: l}
}

func (l *JsNodeList) Item(index int) JsNode {
	jsValue := l.jsValue.Call("item", index)
	return NewNonspecificJsNode(jsValue)
}
