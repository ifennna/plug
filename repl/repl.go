package repl

import (
	"bufio"
	"fmt"
	"io"
	"plug/lexer"
	"plug/parser"
)

const PROMPT = "~> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		if scanned := scanner.Scan(); !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		p := parser.New(lex)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		_, _ = io.WriteString(out, program.String())
		_, _ = io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, "Parser errors:\n")
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
