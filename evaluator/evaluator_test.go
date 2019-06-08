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
		{"-7", -7},
		{"-7", -7},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		testBoolObject(t, testCase.expected, evaluated)
	}
}

func TestBangOperator(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
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
