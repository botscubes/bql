package app

import (
	"github.com/botscubes/bql/internal/token"

	"github.com/botscubes/bql/internal/lexer"
	"go.uber.org/zap"
)

var (
	// input = `x = 2 +X _ 3_123 - 7 y if aelse a IF ELSE TRUE FALSE <= >= == 1 != ! -+ -5`
	input = `x = 2 + 3 
" _ 3123 - 7
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
`
)

func Start(log *zap.SugaredLogger) {
	l := lexer.New(input)

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		log.Debugf("%+v", tok)
	}
}
