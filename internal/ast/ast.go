package ast

import (
	"bytes"

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

	scs := spew.ConfigState{Indent: "  ", DisablePointerAddresses: true}

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

	out.WriteString(as.TokenLiteral() + " ")
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

// Expressions
type IntegerLiteral struct {
	Token token.Token // 5 6
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) ToString() string {
	return il.Token.Literal
}

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
