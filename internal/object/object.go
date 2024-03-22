package object

import "fmt"

type ObjectType = string

type Object interface {
	Type() ObjectType
	ToString() string
}

const (
	INTEGER_OBJ = "INTEGER"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) ToString() string { return fmt.Sprintf("%d", i.Value) }
