package main

import (
	"bytes"
	"plug/repl"
	"strings"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")

	// register functions
	registerCallbacks()
	<-c
}

func registerCallbacks() {
	js.Global().Set("run", js.FuncOf(run))
}

func run(this js.Value, args []js.Value) interface{} {
	code := js.Global().Get("document").Call("getElementById", "input").Get("value").String()
	reader := strings.NewReader(code)
	writer := &bytes.Buffer{}

	repl.Start(reader, writer)

	js.Global().Get("document").Call("getElementById", "output").Set("value", writer.String())

	return true
}
