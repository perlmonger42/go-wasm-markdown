package dom

import "syscall/js"

func SetTitle(title string) {
	js.Global().Get("document").Set("title", title)
}
