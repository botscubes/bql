package app

import (
	"os"

	"github.com/botscubes/bql/api"
	"github.com/botscubes/bql/internal/parser"
	"github.com/botscubes/bql/internal/token"

	"github.com/botscubes/bot-components/context"
	"github.com/botscubes/bql/internal/lexer"
	"go.uber.org/zap"
)

func Start(log *zap.SugaredLogger, fileName string) {

	input, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalw("error opening the file", "error:", err)
	}

	// l := lexer.New(string(input))

	// print_ast(log, l)

	// print_tokens(log, l)

	// p := parser.New(l)
	// program := p.ParseProgram()
	// if len(p.Errors()) != 0 {
	// 	for _, e := range p.Errors() {
	// 		log.Errorln(e)
	// 	}
	// 	return
	// }

	// env := object.NewEnv()
	// ev := evaluator.Eval(program, env)
	// if ev != nil {
	// 	log.Debug(ev.ToString())
	// }

	ctx, passVars, err := prepareCtx()
	if err != nil {
		log.Errorw("prepareCtx", "error", err)
	}

	result, err := api.EvalWithCtx(string(input), ctx, &passVars)
	if err != nil {
		log.Errorw("EvalWithCtx", "error", err)
	} else {
		log.Debugf("RESULT TYPE: %T", result)
		log.Debugf("RESULT VALUE: %+v", result)
	}

	log.Info("Done")
}

func print_tokens(log *zap.SugaredLogger, l *lexer.Lexer) {
	for tok, pos := l.NextToken(); tok.Type != token.EOF; tok, pos = l.NextToken() {
		log.Debugf("Token: %q \t Value: %q \t Pos: %d:%d", tok.Type, tok.Literal, pos.Line, pos.Offset)
	}
}

func print_ast(log *zap.SugaredLogger, l *lexer.Lexer) {
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
}

func prepareCtx() (*context.Context, []string, error) {
	passVars := []string{"x", "s", "m", "b", "a"}
	ctxJson := `
{
	"x": 10,
	"s": "qwerty",
	"b": true,
	"m": {
		"a": 1,
		"b": 2,
		"c": "str",
		"d": true,
		"e": {
			"l2a": 1,
			"l2b": 2
		}
	},
	"a": [1, 2, 3, [true, false], [{
		"s": "qq"
	}]]
}`

	ctx, err := context.NewContextFromJSON([]byte(ctxJson))
	if err != nil {
		return nil, nil, err
	}

	return ctx, passVars, nil
}
