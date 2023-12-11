package lexer

import "github.com/botscubes/bql/internal/token"

type Lexer struct {
	input   string
	ch      byte // current char
	pos     int  // current position (on current char)
	readPos int  // position after current char
	nlsemi  bool // if "true" '\n' translate to ';'
}

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		nlsemi: false,
	}

	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	nlsemi := false

	var tok token.Token
	switch l.ch {
	case '\n':
		tok = newToken(token.SEMICOLON, l.ch)
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			literal := "=="
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.STAR, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			literal := "!="
			tok = token.Token{Type: token.NEQ, Literal: literal}
		} else {
			tok = newToken(token.EXCLAMINATION, l.ch)
		}

	case '%':
		tok = newToken(token.PERCENT, l.ch)
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			literal := "<="
			tok = token.Token{Type: token.LEQ, Literal: literal}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			literal := ">="
			tok = token.Token{Type: token.GEQ, Literal: literal}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAR, l.ch)
	case ')':
		nlsemi = true
		tok = newToken(token.RPAR, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		nlsemi = true
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		nlsemi = true
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		nlsemi = true
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdent()
			tok.Type = token.LookupIdent(tok.Literal)

			if tok.Type == token.IDENT || tok.Type == token.TRUE || tok.Type == token.FALSE {
				l.nlsemi = true
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			l.nlsemi = true
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.nlsemi = nlsemi

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' && !l.nlsemi || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdent() string {
	position := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func (l *Lexer) readNumber() string {
	position := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func (l *Lexer) readString() string {
	position := l.pos + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.pos]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
