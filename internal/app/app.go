package app

import (
	"github.com/botscubes/bql/internal/parser"

	"github.com/botscubes/bql/internal/lexer"
	"go.uber.org/zap"
)

var (
	input = `x = (2+3)*5;
if(x >= 10) {
	t = 1 + 2 * (3 - 125) % 2 / (9 + 1);
} else {
	t = -1;
}

r = add(1, 2 * 8, t, 2 + 3);
`
)

func Start(log *zap.SugaredLogger) {
	l := lexer.New(input)
	p := parser.New(l)

	result := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, e := range p.Errors() {
			log.Errorln(e)
		}
		return
	}

	log.Debug(result.ToString())
	log.Debug(result.Tree())
	log.Info("Done")
}
