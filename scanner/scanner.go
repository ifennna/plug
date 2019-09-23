package scanner

import (
	"bytes"
	"github.com/noculture/plug/evaluator"
	"github.com/noculture/plug/lexer"
	"github.com/noculture/plug/object"
	"github.com/noculture/plug/parser"
	"io"
	"io/ioutil"
)

func Start(in io.Reader, out io.Writer) {
	scanner, _ := ioutil.ReadAll(in)
	env := object.NewEnvironment()

	input := bytes.NewBuffer(scanner).String()

	lex := lexer.New(input)
	p := parser.New(lex)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	evaluator.Out = send(out)

	evaluator.Eval(program, env)
}

func printParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, "Parser errors:\n")
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}

func send(out io.Writer) func(value string) {
	return func (value string) {
		_, _ = io.WriteString(out, value)
	}
}
