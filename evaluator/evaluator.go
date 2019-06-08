package evaluator

import (
	"plug/ast"
	"plug/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		rightExpression := Eval(node.Right)
		return evalPrefixExpression(node.Operator, rightExpression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return referenceBoolObject(node.Value)
	}

	return NULL
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return NULL
	}
}

func evalBangOperator(expression object.Object) object.Object {
	switch expression {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperator(expression object.Object) object.Object {
	if expression.Type() != object.INTEGER {
		return NULL
	}
	value := expression.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func referenceBoolObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
