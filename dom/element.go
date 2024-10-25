package dom

import "syscall/js"

// GoElement models an Element in the virtual DOM, which can be instantiated
// as a real browser DOM element by calling CreateNode().
type GoElement struct {
	tag        string
	attributes []GoAttribute
	children   []GoNode
}

var _ GoNode = &GoElement{}

func NewGoElement(tag string, attrs []GoAttribute, children ...GoNode) *GoElement {
	return &GoElement{tag: tag, attributes: attrs, children: children}
}

func (e *GoElement) AddAttribute(attr GoAttribute) {
	e.attributes = append(e.attributes, attr)
}

func (e *GoElement) AppendChild(child GoNode) {
	e.children = append(e.children, child)
}

// BuildRawDomNode builds and returns a browser DOM element as a js.Value.
func (e *GoElement) BuildRawDomNode() js.Value {
	return BuildRawDomElement(e)
}

// BuildJsNode builds and returns a browser DOM element js.Value wrapped in a JsElement as a generic JsNode.
func (e *GoElement) BuildJsNode() JsNode {
	return e.BuildJsElement()
}

// BuildJsElement builds and returns a browser DOM element js.Value wrapped in a JsElement.
func (e *GoElement) BuildJsElement() *JsElement {
	return NewJsElement(e.BuildRawDomNode())
}

func (e *GoElement) DebugName() string {
	name := "GoElement: <" + e.tag
	for _, attr := range e.attributes {
		if attr.GetName() == "id" {
			name += " id=\"" + attr.DebugValue() + "\""
			break
		}
	}
	return name + ">"
}

func (e *GoElement) DebugJsValue() js.Value {
	// Create a new JavaScript map object containing the element's attributes.
	m := js.Global().Get("Map").New()
	m.Set("tag", e.tag)
	attrs := js.Global().Get("Map").New()
	for _, attr := range e.attributes {
		attrs.Set(attr.GetName(), attr.DebugValue())
	}
	m.Set("attributes", attrs)
	children := js.Global().Get("Array").New()
	for _, child := range e.children {
		children.Call("push", child.DebugJsValue())
	}
	m.Set("children", children)
	return m
}

// JsElement is a strongly typed wrapper around a js.Value representing a DOM element
// (as created by `document.createElement(tagName)` or GoNode.CreateNode()).
type JsElement struct {
	node js.Value
}

var _ JsNode = JsElement{}
var _ JsNodeContainer = JsElement{}

// NewJsElement returns a wrapper of a js.Value referencing an actual browser DOM element.
// Panics if js.Value is not a DOM element (as created by `document.createElement(tagName)`).
func NewJsElement(el js.Value) *JsElement {
	if el.IsNull() || el.IsUndefined() {
		panic("NewJsElement: js.Value is null or undefined")
	} else if el.Get("nodeType").Int() != 1 {
		panic("NewJsElement: js.Value is not a DOM element")
	}
	return &JsElement{node: el}
}

func (e JsElement) isJsNode() {}

func (e JsElement) JsValue() js.Value {
	return e.node
}

func BuildRawDomElement(vElement *GoElement) js.Value {
	jsValueElement := js.Global().Get("document").Call("createElement", vElement.tag)
	for _, val := range vElement.attributes {
		val.SetAttributeOn(jsValueElement)
	}
	fragment := js.Global().Get("document").Call("createDocumentFragment")
	for _, child := range vElement.children {
		fragment.Call("appendChild", child.BuildRawDomNode())
	}
	jsValueElement.Call("appendChild", fragment)
	return jsValueElement
}

func (e JsElement) AppendChild(child GoNode) {
	AppendChild(e, child.BuildJsNode())
}

func (e JsElement) AppendChildren(children ...GoNode) {
	AppendChildren(e, children...)
}

func (e JsElement) QuerySelector(selector string) (node JsNode, ok bool) {
	return QuerySelector(e, selector)
}

func (e JsElement) QuerySelectorAll(selector string) *JsNodeList {
	return QuerySelectorAll(e, selector)
}

func (e JsElement) AppendGoNodes(children ...GoNode) {
	for _, child := range children {
		e.node.Call("appendChild", child.BuildRawDomNode())
	}
}

func (e JsElement) ReplaceChildren(children ...GoNode) {
	e.node.Call("replaceChildren", BuildUntypedNodes(children)...)
}

func GetElementById(id string) (found *JsElement, ok bool) {
	el := js.Global().Get("document").Call("getElementById", id)
	if el.IsNull() || el.IsUndefined() {
		return nil, false
	}
	return NewJsElement(el), true
}
