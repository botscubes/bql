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
		expected int64
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
		testInteger(t, ev, test.expected)
	}
}

func TestEvalBooleanExpresion(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
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
		{`"Abc" == "abc"`, false},
		{`"Abc" == "Abc"`, true},
		{`"Abc" != "Abc"`, false},
		{`"bbc" > "abc"`, true},
		{`"abc" > "bbc"`, false},
		{`"abc" > "abc"`, false},
		{`"abc" < "bbc"`, true},
		{`"abc" <= "bbc"`, true},
		{`"abc" <= "abc"`, true},
		{`"abc" >= "abc"`, true},
		{`"abc" >= "bbc"`, false},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		testBoolean(t, ev, test.expected, test.input)
	}
}

func TestExclaminationOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!false", false},
		{"!!true", true},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		testBoolean(t, ev, test.expected, test.input)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
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
		intVal, ok := test.expected.(int)
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
		expected int64
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
		testInteger(t, ev, test.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"1; true - false; 2", "unknown operator: BOOLEAN - BOOLEAN"},
		{"1; true + false + true + true; 2", "unknown operator: BOOLEAN + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + 3", "type mismatch: BOOLEAN + INTEGER"},
		{"3 * false", "type mismatch: INTEGER * BOOLEAN"},
		{`"Hello" * 3`, "type mismatch: STRING * INTEGER"},
		{`"Hello" * "Earth"`, "unknown operator: STRING * STRING"},
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
		{"true[1]", "index operator not supported: BOOLEAN"},
		{"123[123]", "index operator not supported: INTEGER"},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		err, ok := ev.(*object.Error)
		if !ok {
			t.Errorf("non error object returned: %T - %+v", ev, ev)
		} else {
			if err.Message != test.expected {
				t.Errorf("wrong error message. got: %q expected: %q", err.Message, test.expected)
			}
		}
	}
}

func TestEvalAssignExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"x = 10; x", 10},
		{"x = 10; x = 15; x", 15},
		{"x = 10 * 2; x", 20},
		{"x = 10 * 2; y = x * 3; y", 60},
		{"x = 10; y = x * 2; z = x + y - 20; z", 10},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		testInteger(t, ev, test.expected)
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"x = fn(x){ x }; x(10);", 10},
		{"x = fn(x, y){ return x * y }; x(10, 9);", 90},
		{"x = fn(x, y, z){ return x * (y - z) }; x(10, 9, 1);", 80},
		{"x = fn(x, y){ return x + y }; x(10, x(x(1, 1), x(3, 5)));", 20},
		{"fn(x){ x }(5)", 5},
		{
			`
c = fn(x){
	fn(y) { x + y }
}

a = c(5);
a(4)
`, 9},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		testInteger(t, ev, test.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello Earth"`, "Hello Earth"},
		{`"Hello" + " " + "Earth"`, "Hello Earth"},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		res, ok := ev.(*object.String)
		if !ok {
			t.Errorf("obj not String got:%+v", ev)
			return
		}

		if res.Value != test.expected {
			t.Errorf("obj wrong value. got: %s expected: %s", res.Value, test.expected)
		}
	}
}

func TestArray(t *testing.T) {
	input := "[1, 2, -33, 5+5, 1 + 2 + 3 + 4 * 5]"

	ev := getEvaluated(input)
	res, ok := ev.(*object.Array)
	if !ok {
		t.Errorf("obj not Array got:%+v", ev)
		return
	}

	if len(res.Elements) != 5 {
		t.Errorf("wrong array length got:%d expected 5", len(res.Elements))
		return
	}

	testInteger(t, res.Elements[0], 1)
	testInteger(t, res.Elements[1], 2)
	testInteger(t, res.Elements[2], -33)
	testInteger(t, res.Elements[3], 10)
	testInteger(t, res.Elements[4], 26)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"[1, 2, 5][0]", 1},
		{"[1, 2, 5][2]", 5},
		{"[1, 2, 5][3]", nil},
		{"[1, 2, 5][-1]", nil},
		{"[1, 2, 5][1+1]", 5},
		{"x = 1; [1, 2, 5][x]", 2},
		{"a = [1, 2, 5]; a[0] + a[1] * a[2]", 11},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		intVal, ok := test.expected.(int)
		if ok {
			testInteger(t, ev, int64(intVal))
		} else {
			testNull(t, ev)
		}
	}
}

func TestHashMap(t *testing.T) {
	input := `x = "v";
	{
		"a": 1,
		"bb": 10*10,
		x: 4,
		"qq"+"ww": 123,
		true: 1,
		false: 0
	}`

	ev := getEvaluated(input)
	res, ok := ev.(*object.HashMap)
	if !ok {
		t.Fatalf("non HashMap returned: %T - %+v", ev, ev)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "a"}).HashKey():    1,
		(&object.String{Value: "bb"}).HashKey():   100,
		(&object.String{Value: "v"}).HashKey():    4,
		(&object.String{Value: "qqww"}).HashKey(): 123,
		TRUE.HashKey():  1,
		FALSE.HashKey(): 0,
	}

	if len(res.Pairs) != len(expected) {
		t.Fatalf("wrong len pairs. got: %d expected:%d", len(res.Pairs), len(expected))
	}

	for ek, ev := range expected {
		p, ok := res.Pairs[ek]
		if !ok {
			t.Errorf("not found pair for Key")
		}

		testInteger(t, p.Value, ev)
	}
}

func TestHashMapIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`{"x": 1}["x"]`, 1},
		{`{"x": 1}["y"]`, nil},
		{`{}["y"]`, nil},
		{`{5: 2}[5]`, 2},
		{`{true: 5}[true]`, 5},
		{`{false: -1}[false]`, -1},
		{`y = "key"; {"key": 0}[y]`, 0},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		intVal, ok := test.expected.(int)
		if ok {
			testInteger(t, ev, int64(intVal))
		} else {
			testNull(t, ev)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input         string
		expected      any
		expectIsError bool
	}{
		{`len("abc")`, 3, false},
		{`len("abc" + "efg")`, 6, false},
		{`len("")`, 0, false},
		{`len(1)`, "type of argument not supported: INTEGER", true},
		{`len("a", "b")`, "wrong number of arguments: 2 want: 1", true},
		{`x = "abc"; len(x)`, 3, false},
		{`len([])`, 0, false},
		{`len([1, 2])`, 2, false},
		{`x = [1, 2, 3]; len(x)`, 3, false},
		{`push([], 4)`, []int{4}, false},
		{`push([1, 2, 3], 4)`, []int{1, 2, 3, 4}, false},
		{`push("a", 4)`, "first argument must be ARRAY, got: STRING", true},
		{`first([])`, nil, false},
		{`first([1])`, 1, false},
		{`first([3, 2, 1])`, 3, false},
		{`first("a")`, "argument must be ARRAY, got: STRING", true},
		{`last([])`, nil, false},
		{`last([1])`, 1, false},
		{`last([3, 2, 1])`, 1, false},
		{`last("a")`, "argument must be ARRAY, got: STRING", true},
		{`intToString("a")`, "argument must be INTEGER, got: STRING", true},
		{`intToString(123)`, "123", false},
		{`intToString(1, 2)`, "wrong number of arguments: 2 want: 1", true},
	}

	for _, test := range tests {
		ev := getEvaluated(test.input)
		switch ex := test.expected.(type) {
		case int:
			testInteger(t, ev, int64(ex))
		case nil:
			testNull(t, ev)
		case string:
			if test.expectIsError {
				res, ok := ev.(*object.Error)
				if !ok {
					t.Errorf("object is not Error: %T - %+v", ev, ev)
					continue
				}

				if res.Message != ex {
					t.Errorf("wrong error message: %s expected: %s", res.Message, ex)
				}
			} else {
				res, ok := ev.(*object.String)
				if !ok {
					t.Errorf("object is not String: %T - %+v", ev, ev)
					continue
				}

				if res.Value != ex {
					t.Errorf("wrong string value: %s expected: %s", res.Value, ex)
				}
			}
		case []int:
			res, ok := ev.(*object.Array)
			if !ok {
				t.Errorf("object is not Array: %T - %+v", ev, ev)
				continue
			}

			if len(res.Elements) != len(ex) {
				t.Errorf("wrong number of elements: %d expected: %d", len(res.Elements), len(ex))
				continue
			}

			for i, el := range ex {
				testInteger(t, res.Elements[i], int64(el))
			}
		}
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
