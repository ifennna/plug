package evaluator

import (
	"plug/lexer"
	"plug/object"
	"plug/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"7", 7},
	}

	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		testIntegerObject(t, testCase.expected, evaluated)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		testBoolObject(t, testCase.expected, evaluated)
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, expected int64, evaluated object.Object) bool {
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Errorf("evaluated isn't a plug integer, got %T (%v),", evaluated, evaluated)
		return false
	}
	if result.Value != expected {
		t.Errorf("got wrong integer value, expected %d, got %d", expected, result.Value)
	}

	return true
}

func testBoolObject(t *testing.T, expected bool, evaluated object.Object) bool {
	result, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Errorf("evaluated isn't a plug boolean, got %T (%v),", evaluated, evaluated)
		return false
	}
	if result.Value != expected {
		t.Errorf("got wrong boolean value, expected %t, got %t", expected, result.Value)
	}

	return true
}
