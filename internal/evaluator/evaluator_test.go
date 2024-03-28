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
		{"-5", -5},
		{"-15", -15},
		{"2 + 2", 4},
		{"4 - 2", 2},
		{"2 - 2", 0},
		{"1 + 2 + 3 + 4 + 5 - 1 - 2 - 3", 9},
		{"2 * 3 * 4 * 5 * 6 * 7 * 8 * 9", 362880},
		{"10 + 10 * 2", 30},
		{"(10 + 10) * 2", 40},
		{"100 / 2 * 2 + 5", 105},
		{"100 / (2 * 2) - 200", -175},
		{"5 % 2", 1},
		{"4 % 2", 0},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		testInteger(t, ev, test.excepted)
	}
}

func TestEvalBooleanExpresion(t *testing.T) {
	tests := []struct {
		input    string
		excepted bool
	}{
		{"true", true},
		{"false", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 <= 2", true},
		{"1 <= 1", true},
		{"1 >= 2", false},
		{"1 >= 1", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"(true == false) == false", true},
		{"(1 == 1) == true", true},
		{"(2 > 1) == true", true},
		{"(2 < 1) == false", true},
		{"(1 <= 1) == false", false},
		{"true || false", true},
		{"true && false", false},
		{"true && true", true},
		{"false && false", false},
		{"false || false", false},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		testBoolean(t, ev, test.excepted, test.input)
	}
}

func TestExclaminationOperator(t *testing.T) {
	tests := []struct {
		input    string
		excepted bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!false", false},
		{"!!true", true},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		testBoolean(t, ev, test.excepted, test.input)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		excepted any
	}{
		{"if (true) { 50 }", 50},
		{"if (false) { 50 }", nil},
		{"if (!false) { 50 }", 50},
		{"if (1 == 1) { 50 }", 50},
		{"if (1 > 2) { 50 }", nil},
		{"if (1 < 2) { 50 }", 50},
		{"if (1 < 2) { 50 } else { 100 }", 50},
		{"if (1 > 2) { 50 } else { 100 }", 100},
		{"if (true || false) { 50 } else { 100 }", 50},
		{"if (true && false) { 50 } else { 100 }", 100},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		intVal, ok := test.excepted.(int)
		if ok {
			testInteger(t, ev, int64(intVal))
		} else {
			testNull(t, ev)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		excepted int64
	}{
		{"return 15", 15},
		{"return 25; 1;", 25},
		{"if (true) { return 99 }", 99},
		{`
if (1 > 0) {
	if (2 > 0) {
		return 3
	}

	return 1;
}`, 3},
		{`
if (0 == 0) {
	if (2 == 0) {
		return 3
	}

	return 1;
}`, 1},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		testInteger(t, ev, test.excepted)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		excepted string
	}{
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"1; true - false; 2", "unknown operator: BOOLEAN - BOOLEAN"},
		{"1; true + false + true + true; 2", "unknown operator: BOOLEAN + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + 3", "type mismatch: BOOLEAN + INTEGER"},
		{"3 * false", "type mismatch: INTEGER * BOOLEAN"},
		{"if (3) { 1 }", "non boolean condition in if statement"},
		{`
if (1 > 0) {
	if (2 > 0) {
		return true + false
	}

	return 1;
}`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"x = 10; q", "identifier not found: q"},
		{"ijk", "identifier not found: ijk"},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		err, ok := ev.(*object.Error)
		if !ok {
			t.Errorf("non error object returned: %T - %+v", ev, ev)
		} else {
			if err.Message != test.excepted {
				t.Errorf("wrong error message. got: %q expected: %q", err.Message, test.excepted)
			}
		}
	}
}

func TestEvalAssignExpression(t *testing.T) {
	tests := []struct {
		input    string
		excepted int64
	}{
		{"x = 10; x", 10},
		{"x = 10; x = 15; x", 15},
		{"x = 10 * 2; x", 20},
		{"x = 10 * 2; y = x * 3; y", 60},
		{"x = 10; y = x * 2; z = x + y - 20; z", 10},
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
	env := object.NewEnv()
	return Eval(program, env)
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

func testBoolean(t *testing.T, obj object.Object, expected bool, in string) {
	res, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("obj not Boolean got:%+v", obj)
		return
	}

	if res.Value != expected {
		t.Errorf("obj wrong value. got: %t expected: %t in test: %s", res.Value, expected, in)
	}
}

func testNull(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("obj is not NULL. got: %+v", obj)
	}

	return true
}
