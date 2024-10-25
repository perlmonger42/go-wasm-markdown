package main

import (
	"bytes"
	"syscall/js"

	"github.com/perlmonger42/go-wasm-markdown/dom"
	"github.com/yuin/goldmark"
)

type MarkdownEditor struct {
	Input string
}

func NewMarkdownEditor(input string) *MarkdownEditor {
	return &MarkdownEditor{input}
}

func (m *MarkdownEditor) Render() dom.GoNode {
	return Fragment("markdownEditor",
		Div(
			Attrs(
				Attr("id", "markdownInputDiv"),
				Attr("style", "float: right"),
			),
			TextArea(
				Attrs(
					Attr("id", "markdownInput"),
					Attr("font-family", "monospace"),
					Attr("rows", "14"),
					Attr("cols", "70"),
					AttrFn("oninput", "handleInputChange", handleInputChange),
				),
				Text(m.Input),
			),
		),
		Div(
			Attrs(Attr("id", "markdownRendered")),
			Html(HtmlFromMarkdown(m.Input)),
		),
	)
}

func HtmlFromMarkdown(markdown string) string {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		panic(err)
	}
	// The goldmark README says:
	// "By default, goldmark does not render raw HTML or potentially dangerous links. "
	// So, it should be ok without sanitizing.
	return buf.String()
}

func handleInputChange(this js.Value, _ []js.Value) interface{} {
	inputValue := this.Get("value").String()
	if jsDom, ok := dom.GetElementById("markdownRendered"); ok {
		jsDom.ReplaceChildren(Html(HtmlFromMarkdown(inputValue)))
	}
	return nil
}
