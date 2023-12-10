package app

import (
	"os"

	"github.com/botscubes/bql/internal/parser"

	"github.com/botscubes/bql/internal/lexer"
	"go.uber.org/zap"
)

func Start(log *zap.SugaredLogger, fileName string) {
	input, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalw("error opening the file", "error:", err)
	}

	l := lexer.New(string(input))
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
