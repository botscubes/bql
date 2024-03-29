package ast

import (
	"bytes"
	"strings"

	"github.com/botscubes/bql/internal/token"
	"github.com/davecgh/go-spew/spew"
)

var nl = "\n"

type Node interface {
	TokenLiteral() string
	ToString() string
}

// All statement nodes implement
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) ToString() string {
	// TODO: replace to:
	// r := strings.NewReader("foobar")
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.ToString())
	}

	return out.String()
}

// return pretty AST nodes
func (p *Program) Tree() string {
	var out bytes.Buffer

	scs := spew.ConfigState{Indent: "|   ", DisablePointerAddresses: true}

	for _, s := range p.Statements {
		out.WriteString(nl + scs.Sdump(s) + nl)
	}

	return out.String()
}

// Statements
type AssignStatement struct {
	Name  *Ident
	Value Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return "" }
func (as *AssignStatement) ToString() string {
	var out bytes.Buffer

	out.WriteString(as.Name.TokenLiteral())
	out.WriteString(" = ")

	if as.Value != nil {
		out.WriteString(as.Value.ToString())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) ToString() string {
	if es.Expression != nil {
		return es.Expression.ToString()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token // {
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) ToString() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.ToString())
	}

	return out.String()
}

type ReturnStatement struct {
	Token token.Token // return
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) ToString() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.ToString())
	}

	out.WriteString(";")

	return out.String()
}

// Expressions
type IntegerLiteral struct {
	Token token.Token // 5 6
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) ToString() string     { return il.Token.Literal }

type Boolean struct {
	Token token.Token // TRUE, FALSE
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) ToString() string     { return b.Token.Literal }

type Ident struct {
	Token token.Token // IDENT
	Value string
}

func (i *Ident) expressionNode()      {}
func (i *Ident) TokenLiteral() string { return i.Token.Literal }
func (i *Ident) ToString() string     { return i.Value }

type InfixExpression struct {
	Token    token.Token // +, -, etc
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.ToString())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.ToString())
	out.WriteString(")")

	return out.String()
}

type PrefixExpression struct {
	Token    token.Token // !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.ToString())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token // if
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("if ( ")
	out.WriteString(ie.Condition.ToString())
	out.WriteString(" ) { ")
	out.WriteString(ie.Consequence.ToString())
	out.WriteString(" } ")

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(" { ")
		out.WriteString(ie.Alternative.ToString())
		out.WriteString(" } ")
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // 'fn'
	Parameters []*Ident
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) ToString() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.ToString())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.ToString())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // '('
	FnName    Expression  // Ident
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) ToString() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.ToString())
	}

	out.WriteString(ce.FnName.ToString())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) ToString() string     { return `"` + sl.Token.Literal + `"` }

type ArrayLiteral struct {
	Token    token.Token // '['
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) ToString() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.ToString())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // [
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.ToString())
	out.WriteString("[")
	out.WriteString(ie.Index.ToString())
	out.WriteString("])")

	return out.String()
}
