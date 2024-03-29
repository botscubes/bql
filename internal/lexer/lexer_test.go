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
2 + 3; y = 4
"ABC" "qqq"
if (true) { 1 } else { 0 }
fn(){}
fn(x){ x }

q = fn(x,y,z){
	r = x+y
	return r * z
}

a && true;
b || true;

qwe123;
_ewq;
_99;
{
	"abc": 123
};
`

	tests := []ExpectedToken{
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "2"},
		{token.PLUS, "+"},
		{token.INT, "3"},
		{token.SEMICOLON, "\n"},
		{token.IDENT, "_"},
		{token.INT, "3123"},
		{token.MINUS, "-"},
		{token.INT, "7"},
		{token.SEMICOLON, "\n"},
		{token.IF, "if"},
		{token.IDENT, "aelse"},
		{token.SEMICOLON, "\n"},
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
		{token.SEMICOLON, "\n"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.INT, "3"},
		{token.LEQ, "<="},
		{token.INT, "2"},
		{token.SEMICOLON, "\n"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, "\n"},
		{token.INT, "1"},
		{token.NEQ, "!="},
		{token.INT, "2"},
		{token.SEMICOLON, "\n"},
		{token.INT, "9"},
		{token.GT, ">"},
		{token.INT, "8"},
		{token.SEMICOLON, "\n"},
		{token.INT, "1"},
		{token.LT, "<"},
		{token.INT, "5"},
		{token.SEMICOLON, "\n"},
		{token.EXCLAMINATION, "!"},
		{token.TRUE, "true"},
		{token.NEQ, "!="},
		{token.FALSE, "false"},
		{token.SEMICOLON, "\n"},
		{token.INT, "5"},
		{token.PERCENT, "%"},
		{token.INT, "1"},
		{token.SEMICOLON, "\n"},
		{token.INT, "0"},
		{token.SLASH, "/"},
		{token.INT, "1"},
		{token.SEMICOLON, "\n"},
		{token.STRING, "abc"},
		{token.SEMICOLON, "\n"},
		{token.STRING, "a 1 -2 yy"},
		{token.SEMICOLON, "\n"},
		{token.INT, "2"},
		{token.PLUS, "+"},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
		{token.SEMICOLON, "\n"},
		{token.STRING, "ABC"},
		{token.STRING, "qqq"},
		{token.SEMICOLON, "\n"},
		{token.IF, "if"},
		{token.LPAR, "("},
		{token.TRUE, "true"},
		{token.RPAR, ")"},
		{token.LBRACE, "{"},
		{token.INT, "1"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.INT, "0"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, "\n"},
		{token.FUNC, "fn"},
		{token.LPAR, "("},
		{token.RPAR, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, "\n"},
		{token.FUNC, "fn"},
		{token.LPAR, "("},
		{token.IDENT, "x"},
		{token.RPAR, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, "\n"},
		{token.IDENT, "q"},
		{token.ASSIGN, "="},
		{token.FUNC, "fn"},
		{token.LPAR, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.COMMA, ","},
		{token.IDENT, "z"},
		{token.RPAR, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "r"},
		{token.ASSIGN, "="},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, "\n"},
		{token.RETURN, "return"},
		{token.IDENT, "r"},
		{token.STAR, "*"},
		{token.IDENT, "z"},
		{token.SEMICOLON, "\n"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, "\n"},
		{token.IDENT, "a"},
		{token.LAND, "&&"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "b"},
		{token.LOR, "||"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "qwe123"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "_ewq"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "_99"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "abc"},
		{token.COLON, ":"},
		{token.INT, "123"},
		{token.SEMICOLON, "\n"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, test := range tests {
		tok, _ := l.NextToken()

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
