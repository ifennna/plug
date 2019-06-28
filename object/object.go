package object

import (
	"bytes"
	"fmt"
	"plug/ast"
	"strings"
)

type Type string

const (
	INTEGER      = "INTEGER"
	BOOLEAN      = "BOOLEAN"
	NULL         = "NULL"
	STRING       = "STRING"
	FUNCTION     = "FUNCTION"
	RETURN_VALUE = "RETURN_VALUE"
	ARRAY        = "ARRAY"
	ERROR        = "ERROR"
	BUILTIN      = "BUILTIN"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (fn *Function) Type() Type { return FUNCTION }
func (fn *Function) Inspect() string {
	var out bytes.Buffer
	var params []string

	for _, param := range fn.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") \n")
	out.WriteString(fn.Body.String())
	out.WriteString("\n")

	return out.String()
}

type ReturnValue struct {
	Value Object
}

func (rValue *ReturnValue) Type() Type      { return RETURN_VALUE }
func (rValue *ReturnValue) Inspect() string { return rValue.Value.Inspect() }

type Integer struct {
	Value int64
}

func (int *Integer) Type() Type      { return INTEGER }
func (int *Integer) Inspect() string { return fmt.Sprintf("%d", int.Value) }

type String struct {
	Value string
}

func (s *String) Type() Type      { return STRING }
func (s *String) Inspect() string { return s.Value }

type Boolean struct {
	Value bool
}

func (bool *Boolean) Type() Type      { return BOOLEAN }
func (bool *Boolean) Inspect() string { return fmt.Sprintf("%t", bool.Value) }

type Array struct {
	Elements []Object
}

func (arr *Array) Type() Type { return ARRAY }
func (arr *Array) Inspect() string {
	var out bytes.Buffer
	var elements []string

	for _, element := range arr.Elements {
		elements = append(elements, element.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type Null struct{}

func (null *Null) Type() Type      { return NULL }
func (null *Null) Inspect() string { return fmt.Sprintf("null") }

type Error struct {
	Message string
}

func (e *Error) Type() Type      { return ERROR }
func (e *Error) Inspect() string { return "Error: " + e.Message }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Function BuiltinFunction
}

func (b *Builtin) Type() Type      { return BUILTIN }
func (b *Builtin) Inspect() string { return "builtin function" }
