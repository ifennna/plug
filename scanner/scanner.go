package scanner

import (
	"bytes"
	"io"
	"io/ioutil"
	"plug/evaluator"
	"plug/lexer"
	"plug/object"
	"plug/parser"
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

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		_, _ = io.WriteString(out, evaluated.Inspect())
		_, _ = io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, "Parser errors:\n")
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
