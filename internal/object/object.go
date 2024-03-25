package object

import "fmt"

type ObjectType = string

type Object interface {
	Type() ObjectType
	ToString() string
}

const (
	ERROR_OBJ = "ERROR"
	NULL_OBJ  = "NULL"

	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
)

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) ToString() string { return "Null" }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) ToString() string { return "error: " + e.Message }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) ToString() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) ToString() string { return fmt.Sprintf("%t", b.Value) }
