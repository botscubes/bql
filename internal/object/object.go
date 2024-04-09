package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/botscubes/bql/internal/ast"
)

type ObjectType string

type BuiltinFunction func(args ...Object) Object

type Object interface {
	Type() ObjectType
	ToString() string
}

const (
	ERROR_OBJ = "ERROR"
	NULL_OBJ  = "NULL"

	RETURN_VALUE_OBJ = "RETURN_VALUE"

	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	STRING_OBJ   = "STRING"
	ARRAY_OBJ    = "ARRAY"
	HASH_MAP_OBJ = "HASH_MAP"

	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

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
func (i *Integer) HashKey() HashKey { return HashKey{Type: i.Type(), Value: uint64(i.Value)} }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) ToString() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	if b.Value {
		return HashKey{Type: b.Type(), Value: 1}
	}
	return HashKey{Type: b.Type(), Value: 0}
}

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
	out.WriteString("{ ")
	out.WriteString(f.Body.ToString())
	out.WriteString("} ")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) ToString() string { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) ToString() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.ToString())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type HashMap struct {
	Pairs map[HashKey]HashPair
}

func (h *HashMap) Type() ObjectType { return HASH_MAP_OBJ }
func (h *HashMap) ToString() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, el := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", el.Key.ToString(), el.Value.ToString()))
	}

	out.WriteString("[")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("]")

	return out.String()
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) ToString() string { return "builtin function" }
