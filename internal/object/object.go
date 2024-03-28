package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/botscubes/bql/internal/ast"
)

type ObjectType = string

type Object interface {
	Type() ObjectType
	ToString() string
}

const (
	ERROR_OBJ = "ERROR"
	NULL_OBJ  = "NULL"

	RETURN_VALUE_OBJ = "RETURN_VALUE"

	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"

	FUNCTION_OBJ = "FUNCTION"
)

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) ToString() string { return "Null" }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) ToString() string { return "error: " + e.Message }

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *Return) ToString() string { return r.Value.ToString() }

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

type Function struct {
	Parameters []*ast.Ident
	Body       *ast.BlockStatement
	Env        *Env
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) ToString() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.ToString())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.ToString())

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) ToString() string { return s.Value }
