package app

import (
	"github.com/botscubes/bql/internal/token"

	"github.com/botscubes/bql/internal/lexer"
	"go.uber.org/zap"
)

var (
	input = `x = 2 +X _ 3_123 - 7 y if aelse a IF ELSE TRUE FALSE <= >= == 1 != ! -+ -5`
)

func Start(log *zap.SugaredLogger) {
	l := lexer.New(input)

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		log.Debugf("%+v", tok)
	}
}
