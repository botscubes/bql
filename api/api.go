package api

import (
	"fmt"

	"github.com/botscubes/bot-components/context"
	"github.com/botscubes/bql/internal/evaluator"
	"github.com/botscubes/bql/internal/lexer"
	"github.com/botscubes/bql/internal/object"
	"github.com/botscubes/bql/internal/parser"
)

func EvalWithCtx(code []byte, ctx *context.Context, passVars *[]string) error {
	l := lexer.New(string(code))

	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		var errorslist string
		for _, e := range p.Errors() {
			errorslist += e
		}
		return fmt.Errorf("%s", errorslist)
	}

	env := object.NewEnv()

	env, err := object.ConvertContextToEnv(ctx, env, passVars)
	if err != nil {
		return err
	}

	ev := evaluator.Eval(program, env)
	if ev != nil {
		fmt.Println(ev.ToString())
	}

	return nil
}

// func Eval(code []byte) error {
// 	l := lexer.New(string(code))

// 	p := parser.New(l)
// 	program := p.ParseProgram()
// 	if len(p.Errors()) != 0 {
// 		var errorslist string
// 		for _, e := range p.Errors() {
// 			errorslist += e
// 		}
// 		return fmt.Errorf("%s", errorslist)
// 	}

// 	env := object.NewEnv()
// 	ev := evaluator.Eval(program, env)
// 	if ev != nil {
// 		fmt.Println(ev.ToString())
// 	}

// 	return nil
// }
