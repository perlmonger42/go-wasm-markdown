package dom

import "syscall/js"

type GoTextNode struct {
	text string
}

var _ GoNode = &GoTextNode{}

func NewGoTextNode(text string) *GoTextNode {
	return &GoTextNode{text: text}
}

func (t *GoTextNode) DebugName() string {
	return "GoTextNode: `" + t.text + "`"
}
func (t *GoTextNode) DebugJsValue() js.Value {
	return js.ValueOf("Text: `" + t.text + "`")
}

func (t *GoTextNode) BuildRawDomNode() js.Value {
	return js.Global().Get("document").Call("createTextNode", t.text)
}

func (t *GoTextNode) BuildJsNode() JsNode {
	return t.BuildJsTextNode()
}

func (t *GoTextNode) BuildJsTextNode() JsTextNode {
	return NewJsTextNode(t.BuildRawDomNode())
}

type JsTextNode interface {
	JsNode
	isJsTextNode() // remove this if JsTextNode ever gets anything distinctive
}

type jsTextNode struct {
	node js.Value
}

func NewJsTextNode(textNode js.Value) JsTextNode {
	if textNode.IsNull() || textNode.IsUndefined() {
		panic("NewJsElement: js.Value is null or undefined")
	} else if textNode.Get("nodeType").Int() != 3 {
		panic("NewJsElement: js.Value is not a textNode node")
	}
	return jsTextNode{node: textNode}
}

func (t jsTextNode) JsValue() js.Value {
	return t.node
}

func (t jsTextNode) isJsNode()     {}
func (t jsTextNode) isJsTextNode() {}
