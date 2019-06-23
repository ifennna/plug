package main

import (
	"bytes"
	"plug/repl"
	"strings"
	"syscall/js"
)

func main() {
	println("WASM Go Initialized")

	compile()
}

func compile() {
	code := js.Global().Get("document").Call("getElementById", "input").Get("value").String()
	reader := strings.NewReader(code)
	writer := &bytes.Buffer{}

	repl.Start(reader, writer)

	js.Global().Get("document").Call("getElementById", "output").Set("value", writer.String())
}
