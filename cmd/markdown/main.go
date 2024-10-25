package main

import (
	"github.com/perlmonger42/go-wasm-markdown/dom"
	"syscall/js"
)

func log(args ...interface{}) {
	js.Global().Get("console").Call("log", args...)
}

// RenderPage replaces the contents of the page body.
// I.e., `children[]` replace "..." in the page's `<body>...</body>`.
func RenderPage(children ...dom.GoNode) {
	body := js.Global().Get("document").Call("querySelector", "body")
	body.Call("replaceChildren", dom.BuildUntypedNodes(children)...)
	select {} // run Go forever
}

func logClick(this js.Value, args []js.Value) interface{} {
	if logger, ok := dom.GetElementById("logger"); ok {
		logger.AppendGoNodes(Div(nil, Text("clicked")))
	} else {
		log("no logger div found")
		js.Global().Get("console").Call("log", `<div id="logger"> not found`)
	}
	return nil
}

func setUpPageReloader() {
	url := "ws://" + js.Global().Get("location").Get("host").String() + "/ws"
	log("Connecting to WebSocket at", url)
	ws := js.Global().Get("WebSocket").New(url)
	ws.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		message := args[0].Get("data").String()
		if message == "reload" {
			log("Reloading page...")
			js.Global().Get("location").Call("reload")
		}
		return nil
	}))
}

func main() {
	setUpPageReloader()
	intialInput := "# Markdown Example\n\nThis is a live editor, try editing the Markdown on the right of the page."
	editor := NewMarkdownEditor(intialInput)
	dom.SetTitle("Markdown Demo")
	RenderPage(
		editor.Render(),
		//Div(
		//	Attrs(
		//		Attr("align", "center"),
		//		Attr("id", "logger"),
		//	),
		//),
	)

	select {}
}
