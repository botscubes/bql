package lexer

import (
	"testing"

	"github.com/botscubes/bql/internal/token"
)

type ExpectedToken struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestNextToken(t *testing.T) {
	input := `x = 2 + 3 
_ 3123 - 7
if aelse
if (x == 1) {
	[ 3, 4]
} else {
	3 <= 2
}

1 != 2
9 > 8
1 < 5

!true != false

5 % 1
0/1
"abc"
"a 1 -2 yy"
`

	tests := []ExpectedToken{
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "2"},
		{token.PLUS, "+"},
		{token.INT, "3"},
		{token.IDENT, "_"},
		{token.INT, "3123"},
		{token.MINUS, "-"},
		{token.INT, "7"},
		{token.IF, "if"},
		{token.IDENT, "aelse"},
		{token.IF, "if"},
		{token.LPAR, "("},
		{token.IDENT, "x"},
		{token.EQ, "=="},
		{token.INT, "1"},
		{token.RPAR, ")"},
		{token.LBRACE, "{"},
		{token.LBRACKET, "["},
		{token.INT, "3"},
		{token.COMMA, ","},
		{token.INT, "4"},
		{token.RBRACKET, "]"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.INT, "3"},
		{token.LEQ, "<="},
		{token.INT, "2"},
		{token.RBRACE, "}"},
		{token.INT, "1"},
		{token.NEQ, "!="},
		{token.INT, "2"},
		{token.INT, "9"},
		{token.GT, ">"},
		{token.INT, "8"},
		{token.INT, "1"},
		{token.LT, "<"},
		{token.INT, "5"},
		{token.EXCLAMINATION, "!"},
		{token.TRUE, "true"},
		{token.NEQ, "!="},
		{token.FALSE, "false"},
		{token.INT, "5"},
		{token.PERCENT, "%"},
		{token.INT, "1"},
		{token.INT, "0"},
		{token.SLASH, "/"},
		{token.INT, "1"},
		{token.STRING, "abc"},
		{token.STRING, "a 1 -2 yy"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, test := range tests {
		tok := l.NextToken()

		if tok.Type != test.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong: expected=%q, got=%q",
				i, test.expectedType, tok.Type)
		}

		if tok.Literal != test.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong: expected=%q, got=%q",
				i, test.expectedLiteral, tok.Literal)
		}
	}
}
