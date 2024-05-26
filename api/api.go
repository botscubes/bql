package api

import (
	"fmt"

	"github.com/botscubes/bot-components/context"
	"github.com/botscubes/bql/internal/evaluator"
	"github.com/botscubes/bql/internal/lexer"
	"github.com/botscubes/bql/internal/object"
	"github.com/botscubes/bql/internal/parser"
)

// code - код
// ctx  - контекст
//
// passVars - названия переменных из контекста, которые будут использоваться в коде
// например
// code:
// y = 2
// x + y
//
// переменная x в коде не объявлена, поэтому для успешного выполнения кода, одна должна быть в контексте и в массиве passVars = ["x"]
//
// пример работы есть в файле internal/app/app.go (func prepareCtx()), код в input.txt (запуск: make start)
//
// результат - значение нативного типа Golang и ошибка, если она есть. Если ошибки нет - nil

func EvalWithCtx(code string, ctx *context.Context, passVars *[]string) (any, error) {
	l := lexer.New(code)

	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		var errorslist string
		for _, e := range p.Errors() {
			errorslist += e
		}
		return nil, fmt.Errorf("%s", errorslist)
	}

	env := object.NewEnv()

	env, err := object.ConvertContextToEnv(ctx, env, passVars)
	if err != nil {
		return nil, err
	}

	ev := evaluator.Eval(program, env)
	if ev != nil {
		v, ok := object.ExtractRawValueFromObject(ev)
		if !ok {
			return nil, fmt.Errorf("%s", v)
		}

		return v, nil
	}

	return nil, fmt.Errorf("eval return null")
}
