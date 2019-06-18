package object

import "fmt"

type Type string

const (
	INTEGER             = "INTEGER"
	BOOLEAN             = "BOOLEAN"
	NULL                = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
)

type Object interface {
	Type() Type
	Inspect() string
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

type Boolean struct {
	Value bool
}

func (bool *Boolean) Type() Type      { return BOOLEAN }
func (bool *Boolean) Inspect() string { return fmt.Sprintf("%t", bool.Value) }

type Null struct{}

func (null *Null) Type() Type      { return NULL }
func (null *Null) Inspect() string { return fmt.Sprintf("null") }
