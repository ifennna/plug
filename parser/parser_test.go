package parser

import (
	"plug/ast"
	lexerPackage "plug/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foo 83.338;
`

	lexer := lexerPackage.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program contains %d statements, not 3", len(program.Statements))
	}

	tests := []struct{
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for index, value := range tests {
		statement := program.Statements[index]
		if !testLetStatement(t, statement, value.expectedIdentifier) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, message := range errors {
		t.Errorf("Parser error: %q", message)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement not *ast.LetStatement. got=%T", statement)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}