package parser

import (
	"fmt"
	"strconv"

	"github.com/botscubes/bql/internal/ast"

	"github.com/botscubes/bql/internal/lexer"
	"github.com/botscubes/bql/internal/token"
)

const (
	_ int = iota
	LOWEST
	LOR         // ||
	LAND        // &&
	EQUALS      // ==
	LESSGREATER // > or < or <= or >=
	SUM         // +
	PRODUCT     // * % /
	PREFIX      // -x or !x
	CALL        // call(x) or ( expr )
	INDEX       // [
)

// TODO: create switch and move to token.go
var precedences = map[token.TokenType]int{
	token.LOR:      LOR,
	token.LAND:     LAND,
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.GEQ:      LESSGREATER,
	token.LEQ:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.STAR:     PRODUCT,
	token.PERCENT:  PRODUCT,
	token.LPAR:     CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	peekPos   token.Pos
	curPos    token.Pos

	prefixParsers map[token.TokenType]prefixParseFn
	infixParsers  map[token.TokenType]infixParseFn
	errors        []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	// prefix parse functions
	p.prefixParsers = make(map[token.TokenType]prefixParseFn)
	p.prefixParsers[token.IDENT] = p.parseIdent
	p.prefixParsers[token.INT] = p.parseInteger
	p.prefixParsers[token.MINUS] = p.parsePrefixExpression
	p.prefixParsers[token.EXCLAMINATION] = p.parsePrefixExpression
	p.prefixParsers[token.TRUE] = p.parseBoolean
	p.prefixParsers[token.FALSE] = p.parseBoolean
	p.prefixParsers[token.LPAR] = p.parseGroupedExpression
	p.prefixParsers[token.IF] = p.parseIfExpression
	p.prefixParsers[token.STRING] = p.parseString
	p.prefixParsers[token.LBRACKET] = p.parseArray
	p.prefixParsers[token.FUNC] = p.parseFunction

	// infix parse functions
	p.infixParsers = make(map[token.TokenType]infixParseFn)
	p.infixParsers[token.PLUS] = p.parseInfixExpression
	p.infixParsers[token.MINUS] = p.parseInfixExpression
	p.infixParsers[token.STAR] = p.parseInfixExpression
	p.infixParsers[token.SLASH] = p.parseInfixExpression
	p.infixParsers[token.PERCENT] = p.parseInfixExpression
	p.infixParsers[token.EQ] = p.parseInfixExpression
	p.infixParsers[token.NEQ] = p.parseInfixExpression
	p.infixParsers[token.LEQ] = p.parseInfixExpression
	p.infixParsers[token.GEQ] = p.parseInfixExpression
	p.infixParsers[token.LT] = p.parseInfixExpression
	p.infixParsers[token.GT] = p.parseInfixExpression
	p.infixParsers[token.LOR] = p.parseInfixExpression
	p.infixParsers[token.LAND] = p.parseInfixExpression
	p.infixParsers[token.LPAR] = p.parseCallExpression

	// read curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) newError(e string) {
	mes := fmt.Sprintf("pos: %d:%d: %s", p.curPos.Line, p.curPos.Offset, e)
	p.errors = append(p.errors, mes)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.curPos = p.peekPos
	p.peekToken, p.peekPos = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.newError(fmt.Sprintf("expected next token: %s, got %s", t, p.peekToken.Type))
		return false
	}
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) expectSemi() bool {
	if !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if !p.curTokenIs(token.SEMICOLON) {
			p.newError(fmt.Sprintf("expected ; at end of statement, got %s", p.curToken.Literal))
			return false
		}
		p.nextToken()
	}
	return true
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()

		if ok := p.expectSemi(); !ok {
			break
		}

	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignStatement()
		}

		return p.parseExpressionStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{
		Name: &ast.Ident{Token: p.curToken, Value: p.curToken.Literal},
	}

	// skip ident and =
	p.nextToken()
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()

		if ok := p.expectSemi(); !ok {
			break
		}
	}

	return block
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParsers[p.curToken.Type]
	if prefix == nil {
		p.newError(fmt.Sprintf("prefix parse function for %s not found", p.curToken.Type))
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParsers[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseFunction() ast.Expression {
	node := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAR) {
		return nil
	}

	node.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	node.Body = p.parseBlockStatement()

	return node
}

func (p *Parser) parseFunctionParameters() []*ast.Ident {
	identifiers := []*ast.Ident{}

	if p.peekTokenIs(token.RPAR) {
		p.nextToken()
		return identifiers
	}

	for {
		p.nextToken()
		if !p.curTokenIs(token.IDENT) {
			p.newError(fmt.Sprintf("failed parse %q as ident", p.curToken.Literal))
			return nil
		}

		ident := &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
		if !p.peekTokenIs(token.COMMA) {
			break
		}

		p.nextToken()
	}

	if !p.expectPeek(token.RPAR) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseInteger() ast.Expression {
	node := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.newError(fmt.Sprintf("failed parse %q as integer", p.curToken.Literal))
		return nil
	}

	node.Value = value

	return node
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseIdent() ast.Expression {
	return &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAR) {
		return nil
	}

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	prec := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(prec)

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAR) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAR) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, FnName: fn}
	exp.Arguments = p.parseExpressionList(token.RPAR)
	return exp
}

func (p *Parser) parseExpressionList(endToken token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(endToken) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(endToken) {
		return nil
	}

	return list
}

func (p *Parser) parseString() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArray() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}
