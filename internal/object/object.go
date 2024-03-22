package object

import "fmt"

type ObjectType = string

type Object interface {
	Type() ObjectType
	ToString() string
}

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) ToString() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return INTEGER_OBJ }
func (b *Boolean) ToString() string { return fmt.Sprintf("%t", b.Value) }
