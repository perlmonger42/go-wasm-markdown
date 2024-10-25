package main

import "github.com/perlmonger42/go-wasm-markdown/dom"

func Attr(name, value string) dom.GoAttribute {
	return dom.NewStringAttribute(name, value)
}

func AttrFn(name string, fnName string, value dom.Callback) dom.GoAttribute {
	return dom.NewFunctionAttribute(name, fnName, value)
}

func Text(text string) *dom.GoTextNode {
	return dom.NewGoTextNode(text)
}

func Attrs(attributes ...dom.GoAttribute) []dom.GoAttribute {
	return attributes
}

func Fragment(debugName string, children ...dom.GoNode) *dom.GoFragment {
	return dom.NewGoFragment(debugName, children...)
}

func Div(attributes []dom.GoAttribute, children ...dom.GoNode) *dom.GoElement {
	return dom.NewGoElement("div", attributes, children...)
}

func Html(html string) dom.GoNode {
	return dom.NewGoHtmlNode(html)
}

func TextArea(attributes []dom.GoAttribute, children ...dom.GoNode) *dom.GoElement {
	return dom.NewGoElement("textarea", attributes, children...)
}
