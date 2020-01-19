// +build js, wasm

package main

import (
	"syscall/js"
)

const (
	loadEventName = "0xmeshtest"
)

func main() {
	triggerLoadEvent()

	select {}
}

// triggerLoadEvent triggers the global load event to indicate that the Wasm is
// done loading.
func triggerLoadEvent() {
	event := js.Global().Get("document").Call("createEvent", "Event")
	event.Call("initEvent", loadEventName, true, true)
	js.Global().Call("dispatchEvent", event)
}
