package evaluator

import (
	"fmt"
	"github.com/noculture/plug/ast"
	"github.com/noculture/plug/object"
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
	case *ast.ForStatement:
		return evalForLoop(node, env)
	case *ast.FunctionLiteral:
		parameters := node.Parameters
		body := node.Body
		return &object.Function{Parameters: parameters, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		arguments := evalExpressions(node.Arguments, env)
		// catch errors
		if len(arguments) == 1 && isError(arguments[0]) {
			return arguments[0]
		}
		return applyFunction(function, arguments)
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
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		// catch errors
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
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
			if resultType == object.RETURN_VALUE || resultType == object.ERROR {
				return result
			}
		}
	}

	return result
}

func evalForLoop(statement *ast.ForStatement, environment *object.Environment) object.Object {
	var body object.Object

	loopNumber := statement.Range.Arguments[0].(*ast.IntegerLiteral)
	value := int(loopNumber.Value)
	for i := 0; i < value; i++ {
		environment.Set(statement.Index.Value, &object.Integer{Value: int64(i)})
		body = evalBlockStatement(statement.Body, environment)

		if body != nil {
			returnType := body.Type()
			if returnType == object.RETURN_VALUE || returnType == object.ERROR {
				return body
			}
		}
	}
	return body
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := Eval(expression, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		innerEnv := createFunctionScope(function, args)
		evaluated := Eval(function.Body, innerEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return function.Function(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func createFunctionScope(fn *object.Function, arguments []object.Object) *object.Environment {

	// passing the function's environment allow for closures, we still have the
	// function's bindings ling after it has finished execution
	env := object.NewEnclosedEvironment(fn.Env)

	for paramIndex, param := range fn.Parameters {
		env.Set(param.Value, arguments[paramIndex])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
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

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arr := array.(*object.Array)
	indexValue := index.(*object.Integer).Value
	max := int64(len(arr.Elements) - 1)

	if indexValue < 0 || indexValue > max {
		return NULL
	}

	return arr.Elements[indexValue]
}

// For expressions that resolve to boolean, direct comparison can be carried out since there are only
// two boolean objects. In other cases the values have to be unwrapped and compared instead.
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)
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

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value
	return &object.String{Value: leftValue + rightValue}
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
	if value, ok := env.Get(node.Value); ok {
		return value
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

func referenceBoolObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
