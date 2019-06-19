package evaluator

import (
	"fmt"
	"plug/ast"
	"plug/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	case *ast.LetStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Set(node.Name.Value, value)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.InfixExpression:
		leftExpression := Eval(node.Left, env)
		if isError(leftExpression) {
			return leftExpression
		}
		rightExpression := Eval(node.Right, env)
		if isError(rightExpression) {
			return rightExpression
		}
		return evalInfixExpression(node.Operator, leftExpression, rightExpression)
	case *ast.PrefixExpression:
		rightExpression := Eval(node.Right, env)
		if isError(rightExpression) {
			return rightExpression
		}
		return evalPrefixExpression(node.Operator, rightExpression)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return referenceBoolObject(node.Value)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		// if we have an error or a return statement, return the value and
		// ignore the rest of the code within the scope
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			resultType := result.Type()
			if resultType == object.RETURN_VALUE_OBJECT || resultType == object.ERROR_OBJECT {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ifExp *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ifExp.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ifExp.Consequence, env)
	} else if ifExp.Alternative != nil {
		return Eval(ifExp.Alternative, env)
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
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
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
		return newError("unknown operator: -%s", expression.Type())
	}
	value := expression.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	value, ok := env.Get(node.Value)
	if !ok {
		return newError("variable has not been declared: " + node.Value)
	}
	return value
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJECT
	}
	return false
}

func referenceBoolObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
