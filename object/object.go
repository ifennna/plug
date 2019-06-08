package object

import "fmt"

type Type string

const (
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	NULL    = "NULL"
)

type Object interface {
	Type() Type
	Inspect() string
}

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
