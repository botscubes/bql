package evaluator

import (
	"testing"

	"github.com/botscubes/bql/internal/lexer"
	"github.com/botscubes/bql/internal/object"
	"github.com/botscubes/bql/internal/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		excepted int64
	}{
		{"4", 4},
		{"12", 12},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		testInteger(t, ev, test.excepted)
	}
}

func getEvaluated(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testInteger(t *testing.T, obj object.Object, expected int64) {
	res, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj not Integer got:%+v", obj)
		return
	}

	if res.Value != expected {
		t.Errorf("obj wrong value. got: %d expected: %d", res.Value, expected)
	}

}
