package parser

import (
	"fmt"
	"plug/ast"
	lexerPackage "plug/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	testCases := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foo = y;", "foo", "y"},
	}

	for _, testCase := range testCases {
		program := setup(testCase.input, t)

		if len(program.Statements) != 1 {
			t.Fatalf("program contains %d statements, not 1", len(program.Statements))
		}
		statement := program.Statements[0]
		if !testLetStatement(t, statement, testCase.expectedIdentifier) {
			return
		}
		value := statement.(*ast.LetStatement).Value
		if !testLiteralExpression(t, value, testCase.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	testCases := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, testCase := range testCases {
		program := setup(testCase.input, t)

		if len(program.Statements) != 1 {
			t.Fatalf("program contains %d statements, not 1", len(program.Statements))
		}
		statement := program.Statements[0]
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, got %T", statement)
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement token not 'return', got %q", returnStatement.TokenLiteral())
		}
		if !testLiteralExpression(t, returnStatement.ReturnValue, testCase.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	program := setup(input, t)

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

	program := setup(input, t)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements, got %d", len(program.Statements))
	}
	statement := getStatement(program, t)

	if !testIntegerLiteral(t, statement.Expression, 5) {
		return
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	program := setup(input, t)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := statement.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expresssion is not a string literal, got %T", statement.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal value is not %q, got %q", "hello world", literal.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "false;"

	program := setup(input, t)

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
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, testCase := range testCases {
		program := setup(testCase.input, t)

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
		if !testLiteralExpression(t, expression.Right, testCase.value) {
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
		program := setup(testCase.input, t)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements contains %d statements, not %d", len(program.Statements), 1)
		}
		statement := getStatement(program, t)
		if !testInfixExpression(t, statement.Expression, testCase.leftValue, testCase.operator, testCase.rightValue) {
			return
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := setup(input, t)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements contains %d statements, not %d", len(program.Statements), 1)
	}
	statement := getStatement(program, t)
	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expression is not an ast.IfExpression, got %T", statement.Expression)
	}
	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement(s), got %d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence is not an ast.ExpressionStatement, got %T", expression.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if expression.Alternative != nil {
		t.Errorf("expression alternative is not nil, got %+v", expression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program := setup(input, t)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	statement := getStatement(program, t)
	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement expression is not ast.IfExpression. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("expression.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			expression.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func(a, b) {a * b;}`

	program := setup(input, t)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements contains %d statements, not %d", len(program.Statements), 1)
	}
	statement := getStatement(program, t)
	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement expression is not ast.FunctionLiteral. got=%T", statement.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong, epected 2, got %d", len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "a")
	testLiteralExpression(t, function.Parameters[1], "b")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("expected 1 function body statement(s), got %d", len(function.Body.Statements))
	}
	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("body statement not ast.ExpressionStatement, get %T", function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStatement.Expression, "a", "*", "b")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "func() {};", expectedParams: []string{}},
		{input: "func(x) {};", expectedParams: []string{"x"}},
		{input: "func(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, testCase := range tests {
		program := setup(testCase.input, t)
		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(testCase.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(testCase.expectedParams), len(function.Parameters))
		}
		for i, ident := range testCase.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2*3, 4+5)"

	program := setup(input, t)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements contains %d statements, not %d", len(program.Statements), 1)
	}
	statement := getStatement(program, t)
	expression, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("statement expression is not ast.CallExpression. got=%T", statement.Expression)
	}

	if !testIdentifier(t, expression.Function, "add") {
		return
	}
	if len(expression.Arguments) != 3 {
		t.Fatalf("wrong length of arguments, got %d", len(expression.Arguments))
	}

	testLiteralExpression(t, expression.Arguments[0], 1)
	testInfixExpression(t, expression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expression.Arguments[2], 4, "+", 5)
}

func TestArrayLiteralParsing(t *testing.T) {
	input := "[1, 2 * 3, 4 + 5]"
	program := setup(input, t)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := statement.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expression is not an ast.ArrayLiteral. got=%T", statement.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("length of elements is not 3, got %d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 3)
	testInfixExpression(t, array.Elements[2], 4, "+", 5)
}

func TestIndexExpressionParsing(t *testing.T) {
	input := "myArray[1 + 1]"
	program := setup(input, t)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := statement.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("expression is not an ast.IndexExpression. got %T", statement.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, testCase := range tests {
		lexer := lexerPackage.New(testCase.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		actual := program.String()
		if actual != testCase.expected {
			t.Errorf("expected=%q, got=%q", testCase.expected, actual)
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
	boolean, ok := expression.(*ast.Boolean)

	if !ok {
		t.Errorf("expression is not an ast.Boolean, got %T", expression)
		return false
	}
	if boolean.Value != value {
		t.Errorf("boolean value is not %t, got %t", value, boolean.Value)
		return false
	}
	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, boolean.TokenLiteral())
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

func setup(input string, t *testing.T) *ast.Program {
	lexer := lexerPackage.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	return program
}
