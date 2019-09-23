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

	compile()
}

func compile() {
	code := js.Global().Get("document").Call("getElementById", "input").Get("value").String()
	reader := strings.NewReader(code)
	writer := &bytes.Buffer{}

	startTime := time.Now()
	scanner.Start(reader, writer)
	elapsed := time.Since(startTime)

	writer.WriteString("\n Time elapsed: " + elapsed.String())

	js.Global().Get("document").Call("getElementById", "output").Set("value", writer.String())
}
