package evaluator

import (
	"plug/lexer"
	"plug/object"
	"plug/parser"
	"testing"
)

func TestEvaluation(t *testing.T) {
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
