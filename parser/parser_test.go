package parser

import (
	"fmt"
	"plug/ast"
	lexerPackage "plug/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foo = 83.338;
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

	tests := []struct {
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

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	lexer := lexerPackage.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf("program contains %d statements, not 3", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, got %T", statement)
			continue
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement token not 'return', got %q", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	lexer := lexerPackage.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements, got %d", len(program.Statements))
	}
	statement := getStatement(program, t)

	if !testIdentifier(t, statement.Expression, "foobar") {
		return
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lexer := lexerPackage.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements, got %d", len(program.Statements))
	}
	statement := getStatement(program, t)

	if !testIntegerLiteral(t, statement.Expression, 5) {
		return
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "false;"

	lexer := lexerPackage.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements, got %d", len(program.Statements))
	}
	statement := getStatement(program, t)

	if !testBoolean(t, statement.Expression, false) {
		return
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	testCases := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, testCase := range testCases {
		lexer := lexerPackage.New(testCase.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements contains %d statements, not %d", len(program.Statements), 1)
		}
		statement := getStatement(program, t)
		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("statement is not a ast.PrefixExpression, got %T", statement.Expression)
		}
		if expression.Operator != testCase.operator {
			t.Fatalf("expression operator is not an %s, got %s", testCase.operator, expression.Operator)
		}
		if !testLiteralExpression(t, expression.Right, testCase.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	testCases := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, testCase := range testCases {
		lexer := lexerPackage.New(testCase.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements contains %d statements, not %d", len(program.Statements), 1)
		}
		statement := getStatement(program, t)
		if !testInfixExpression(t, statement.Expression, testCase.leftValue, testCase.operator, testCase.rightValue) {
			return
		}
	}
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch value := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(value))
	case int64:
		return testIntegerLiteral(t, expression, value)
	case string:
		return testIdentifier(t, expression, value)
	case bool:
		return testBoolean(t, expression, value)
	}
	t.Errorf("type of expression not handled. got=%T", expression)
	return false
}

func testInfixExpression(t *testing.T, expression ast.Expression, left interface{},
	operator string, right interface{}) bool {

	correctExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("statement not *ast.InfixExpression. got %T(%s)", expression, expression)
	}
	if !testLiteralExpression(t, correctExpression.Left, left) {
		return false
	}
	if correctExpression.Operator != operator {
		t.Fatalf("expression operator is not an %s, got %s", operator, correctExpression.Operator)
		return false
	}
	if !testLiteralExpression(t, correctExpression.Right, right) {
		return false
	}

	return true
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

func testIntegerLiteral(t *testing.T, literal ast.Expression, value int64) bool {
	integer, ok := literal.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("integer literal is not an ast.IntegerLiteral, got %T", literal)
		return false
	}
	if integer.Value != value {
		t.Errorf("integer value is not %d, got %d", value, integer.Value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer token literal is not %d, got %s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)

	if !ok {
		t.Errorf("expression is not an ast.Identifier, got %T", expression)
		return false
	}
	if identifier.Value != value {
		t.Errorf("integer value is not %s, got %s", value, identifier.Value)
		return false
	}
	if identifier.TokenLiteral() != value {
		t.Errorf("integer token literal is not %s, got %s", value, identifier.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, expression ast.Expression, value bool) bool {
	identifier, ok := expression.(*ast.Boolean)

	if !ok {
		t.Errorf("expression is not an ast.Identifier, got %T", expression)
		return false
	}
	if identifier.Value != value {
		t.Errorf("integer value is not %t, got %t", value, identifier.Value)
		return false
	}

	return true
}

func getStatement(program *ast.Program, t *testing.T) *ast.ExpressionStatement {
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement, got %T", program.Statements[0])
	}
	return statement
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
