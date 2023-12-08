package app

import (
	"github.com/botscubes/bql/internal/parser"

	"github.com/botscubes/bql/internal/lexer"
	"go.uber.org/zap"
)

var (
	input = `x = 3+1*2*4+5; y = 2; true-1;
	
if(a == b) {2 / 3+1*2%1} else { y = x - 1 * 3}

(2 + 3) * 6
y < 1

!x != !y
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

	log.Debug(result.String())
}
