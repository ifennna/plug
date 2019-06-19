package object

import (
	"bytes"
	"fmt"
	"plug/ast"
	"strings"
)

type Type string

const (
	INTEGER             = "INTEGER"
	BOOLEAN             = "BOOLEAN"
	NULL                = "NULL"
	STRING_OBJECT       = "STRING"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
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

func (fn *Function) Type() Type { return FUNCTION_OBJECT }
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

func (rValue *ReturnValue) Type() Type      { return RETURN_VALUE_OBJECT }
func (rValue *ReturnValue) Inspect() string { return rValue.Value.Inspect() }

type Integer struct {
	Value int64
}

func (int *Integer) Type() Type      { return INTEGER }
func (int *Integer) Inspect() string { return fmt.Sprintf("%d", int.Value) }

type String struct {
	Value string
}

func (s *String) Type() Type      { return STRING_OBJECT }
func (s *String) Inspect() string { return s.Value }

type Boolean struct {
	Value bool
}

func (bool *Boolean) Type() Type      { return BOOLEAN }
func (bool *Boolean) Inspect() string { return fmt.Sprintf("%t", bool.Value) }

type Null struct{}

func (null *Null) Type() Type      { return NULL }
func (null *Null) Inspect() string { return fmt.Sprintf("null") }

type Error struct {
	Message string
}

func (e *Error) Type() Type      { return ERROR_OBJECT }
func (e *Error) Inspect() string { return "Error: " + e.Message }
