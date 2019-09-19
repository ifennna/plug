package main

import (
	"bytes"
	"github.com/noculture/plug/scanner"
	"strings"
	"syscall/js"
	"time"
)

func main() {
	println("WASM Go Initialized")

	startTime := time.Now()
	compile()
	elapsed := time.Since(startTime)

	println(elapsed.String())
}

func compile() {
	code := js.Global().Get("document").Call("getElementById", "input").Get("value").String()
	reader := strings.NewReader(code)
	writer := &bytes.Buffer{}

	scanner.Start(reader, writer)

	js.Global().Get("document").Call("getElementById", "output").Set("value", writer.String())
}
