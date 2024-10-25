package dom

import "syscall/js"

// GoFragment is a GoNode that represents a GoFragment in the virtual DOM.
type GoFragment struct {
	debugName string
	children  []GoNode
}

var _ GoNode = &GoFragment{}

func NewGoFragment(debugName string, children ...GoNode) *GoFragment {
	return &GoFragment{children: children}
}

func (node *GoFragment) DebugName() string {
	return node.debugName
}

func (node *GoFragment) DebugJsValue() js.Value {
	// Create a new JavaScript map object describing the fragment.
	m := js.Global().Get("Map").New()
	m.Set("fragment", node.debugName)
	children := js.Global().Get("Array").New()
	for _, child := range node.children {
		children.Call("push", child.DebugJsValue())
	}
	m.Set("children", children)
	return m
}

func (node *GoFragment) AppendChild(child GoNode) {
	node.children = append(node.children, child)
}

func (node *GoFragment) BuildRawDomNode() js.Value {
	frag := js.Global().Get("document").Call("createDocumentFragment")
	for _, child := range node.children {
		frag.Call("appendChild", child.BuildRawDomNode())
	}
	return frag
}

func (node *GoFragment) BuildJsNode() JsNode {
	return node.BuildJsFragment()
}

func (node *GoFragment) BuildJsFragment() *JsFragment {
	return NewJsFragment(node.debugName, node.BuildRawDomNode())
}

type JsFragment struct {
	debugName string
	node      js.Value
}

var _ JsNode = &JsFragment{}

func NewJsFragment(debugName string, fragment js.Value) *JsFragment {
	if fragment.IsNull() || fragment.IsUndefined() {
		panic("NewJsFragment: js.Value is null or undefined")
	} else if fragment.Get("nodeType").Int() != 11 {
		panic("NewJsFragment: js.Value is not a fragment node")
	}
	return &JsFragment{debugName: debugName, node: fragment}
}

func (f JsFragment) isJsNode() {}

func (f JsFragment) DebugName() string {
	return f.debugName
}

func (f JsFragment) JsValue() js.Value {
	return f.node
}
