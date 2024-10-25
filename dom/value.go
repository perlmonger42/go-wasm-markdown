package dom

// We may not need this hideMyValue interface and its implementations.

import "syscall/js"

// hideMyValue is a strongly typed wrapper around js.Value.
type hideMyValue interface {
	JsValue() js.Value
}

// hideStringValue is a strongly typed wrapper around a string.
type hideStringValue struct {
	s string
}

func (v hideStringValue) JsValue() js.Value {
	return js.ValueOf(v.s)
}

func NewhideStringValue(s string) hideMyValue {
	return &hideStringValue{s: s}
}

// hideFunctionValue is a strongly typed wrapper around a function.
type hideFunctionValue struct {
	// You cannot define methods on a function type in Go, so we wrap it in a struct.
	fn Callback
}

func (v hideFunctionValue) JsValue() js.Value {
	return js.ValueOf(js.FuncOf(v.fn))
}

func NewhideFunctionValue(fn Callback) hideFunctionValue {
	return hideFunctionValue{fn: fn}
}

// hideElementValue is a strongly typed wrapper around a js.Value representing a DOM element
// (as created by `document.createElement(tagName)`).
type hideElementValue struct {
	el js.Value
}

func (v hideElementValue) JsValue() js.Value {
	return v.el
}

func NewhideElementValue(el js.Value) hideMyValue {
	return hideElementValue{el: el}
}

// hideTextNodeValue is a strongly typed wrapper around a js.Value representing a DOM text node
// (as created by `document.createTextNode("hello")`).
type hideTextNodeValue struct {
	tn js.Value
}

func (v hideTextNodeValue) JsValue() js.Value {
	return v.tn
}

func NewhideTextNodeValue(tn js.Value) hideMyValue {
	return hideTextNodeValue{tn: tn}
}
