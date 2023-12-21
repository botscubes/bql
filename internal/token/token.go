package token

type TokenType = string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"  // x, t, add
	INT    = "INT"    // 123
	STRING = "STRING" // "abcde"

	ASSIGN        = "="
	PLUS          = "+"
	MINUS         = "-"
	STAR          = "*"
	SLASH         = "/"
	EXCLAMINATION = "!"
	PERCENT       = "%"

	EQ  = "=="
	NEQ = "!="
	LEQ = "<="
	GEQ = ">="
	LT  = "<"
	GT  = ">"

	COMMA     = ","
	SEMICOLON = ";"

	LPAR     = "("
	RPAR     = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// keywords
	IF     = "IF"
	ELSE   = "ELSE"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	FUNC   = "FUNCTION"
	RETURN = "RETURN"
)

var keywords = map[string]TokenType{
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"fn":     FUNC,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
