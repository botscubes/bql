package app

import (
	"github.com/botscubes/bql/internal/parser"

	"github.com/botscubes/bql/internal/lexer"
	"go.uber.org/zap"
)

var (
	input = `x = (2+3)*5`
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
