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
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.InfixExpression:
		leftExpression := Eval(node.Left)
		rightExpression := Eval(node.Right)
		return evalInfixExpression(node.Operator, leftExpression, rightExpression)
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

func evalIfExpression(ifExp *ast.IfExpression) object.Object {
	condition := Eval(ifExp.Condition)

	if isTruthy(condition) {
		return Eval(ifExp.Consequence)
	} else if ifExp.Alternative != nil {
		return Eval(ifExp.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

// For expressions that resolve to boolean, direct comparison can be carried out since there are only
// two boolean objects. In other cases the values have to be unwrapped and compared instead.
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return referenceBoolObject(left == right)
	case operator == "!=":
		return referenceBoolObject(left != right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case ">":
		return referenceBoolObject(leftValue > rightValue)
	case "<":
		return referenceBoolObject(leftValue < rightValue)
	case "==":
		return referenceBoolObject(leftValue == rightValue)
	case "!=":
		return referenceBoolObject(leftValue != rightValue)
	default:
		return NULL
	}
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
