package parser

import (
	"strconv"
	"testing"

	"github.com/botscubes/bql/internal/ast"
	"github.com/botscubes/bql/internal/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) != 0 {
		t.Errorf("parser has %d errors:", len(errors))
		for _, e := range errors {
			t.Errorf("parser error: %q", e)
		}
		t.FailNow()
	}
}

func TestParseAssignStatement(t *testing.T) {
	tests := []struct {
		input string
		ident string
		value any
	}{
		{"x = 56", "x", 56},
		{"y = x", "y", "x"},
		{"k = true", "k", true},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		result := p.ParseProgram()
		checkParserErrors(t, p)

		if len(result.Statements) != 1 {
			t.Fatalf("program has incorrect number of statements. got:%d",
				len(result.Statements))
		}

		stmt, ok := result.Statements[0].(*ast.AssignStatement)
		if !ok {
			t.Fatalf("result.Statements[0] is not ast.AssignStatement. got:%T",
				result.Statements[0])
		}

		if !testIdent(t, stmt.Name, test.ident) {
			return
		}

		if !testLiteralExpression(t, stmt.Value, test.value) {
			return
		}

	}
}

func TestParseReturnStatement(t *testing.T) {
	tests := []struct {
		input string
		value any
	}{
		{"return x", "x"},
		{"return true", true},
		{"return 4;", 4},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		result := p.ParseProgram()
		checkParserErrors(t, p)

		if len(result.Statements) != 1 {
			t.Fatalf("program has incorrect number of statements. got:%d",
				len(result.Statements))
		}

		returnStmt, ok := result.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("result.Statements[0] is not ast.ReturnStatement. got:%T",
				result.Statements[0])
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral is not 'return'. got:%s",
				returnStmt.TokenLiteral())
		}

		if !testLiteralExpression(t, returnStmt.Value, test.value) {
			return
		}
	}
}

func TestParsePrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    any
	}{
		{"-2", "-", 2},
		{"!1", "!", 1},
		{"!true", "!", true},
		{"!false", "!", false},
		{"!xyz", "!", "xyz"},
		{"-xyz", "-", "xyz"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		result := p.ParseProgram()
		checkParserErrors(t, p)

		if len(result.Statements) != 1 {
			t.Fatalf("program has incorrect number of statements. got:%d",
				len(result.Statements))
		}

		stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
				result.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got:%T", stmt.Expression)
		}

		if expr.Operator != test.operator {
			t.Fatalf("expr.Operator is not %s. got:%T",
				test.operator, expr.Operator)
		}

		if !testLiteralExpression(t, expr.Right, test.value) {
			return
		}
	}
}

func TestParseInfixExpression(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"2 + 2", 2, "+", 2},
		{"2 - 2", 2, "-", 2},
		{"2 * 2", 2, "*", 2},
		{"2 / 2", 2, "/", 2},
		{"2 % 2", 2, "%", 2},
		{"2 > 2", 2, ">", 2},
		{"2 < 2", 2, "<", 2},
		{"2 == 2", 2, "==", 2},
		{"2 != 2", 2, "!=", 2},
		{"2 >= 2", 2, ">=", 2},
		{"2 <= 2", 2, "<=", 2},
		{"abc + foo", "abc", "+", "foo"},
		{"abc - foo", "abc", "-", "foo"},
		{"abc * foo", "abc", "*", "foo"},
		{"abc / foo", "abc", "/", "foo"},
		{"abc % foo", "abc", "%", "foo"},
		{"abc > foo", "abc", ">", "foo"},
		{"abc < foo", "abc", "<", "foo"},
		{"abc == foo", "abc", "==", "foo"},
		{"abc != foo", "abc", "!=", "foo"},
		{"abc >= foo", "abc", ">=", "foo"},
		{"abc <= foo", "abc", "<=", "foo"},
		{"true == true", true, "==", true},
		{"false != true", false, "!=", true},
		{"false == false", false, "==", false},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		result := p.ParseProgram()
		checkParserErrors(t, p)

		if len(result.Statements) != 1 {
			t.Fatalf("program has incorrect number of statements. got:%d",
				len(result.Statements))
		}

		stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
				result.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, test.leftValue,
			test.operator, test.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"2 + 3",
			"(2 + 3)",
		},
		{
			"-5 + 3",
			"((-5) + 3)",
		},
		{
			"-(2 + 4)",
			"(-(2 + 4))",
		},
		{
			"(1 + 2) + 3",
			"((1 + 2) + 3)",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a % b + c",
			"((a % b) + c)",
		},
		{
			"a % b * c",
			"((a % b) * c)",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"5 <= 4 != (5 == 5)",
			"((5 <= 4) != (5 == 5))",
		},
		{
			"3 + 6 * 5 == 3 * 2 + 4 * 5",
			"((3 + (6 * 5)) == ((3 * 2) + (4 * 5)))",
		},
		{
			"(3 + 6) * 5 == 3 * (2 + 4) * 5",
			"(((3 + 6) * 5) == ((3 * (2 + 4)) * 5))",
		},
		{
			"!true",
			"(!true)",
		},
		{
			"!(true != false)",
			"(!(true != false))",
		},
		{
			"a || b && c",
			"(a || (b && c))",
		},
		{
			"(a || b) && c || d",
			"(((a || b) && c) || d)",
		},
		{
			"a > b && c",
			"((a > b) && c)",
		},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		result := p.ParseProgram()
		checkParserErrors(t, p)

		text := result.ToString()
		if text != test.expected {
			t.Errorf("expected=%q, got:%q", test.expected, text)
		}
	}
}

func TestParseIfExpression(t *testing.T) {
	tests := []struct {
		input           string
		withAlternative bool
	}{
		{"if (x == y) { 1 }", false},
		{"if (x == y) { 1 } else { 0 }", true},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		result := p.ParseProgram()
		checkParserErrors(t, p)

		if len(result.Statements) != 1 {
			t.Fatalf("program has incorrect number of statements. got:%d",
				len(result.Statements))
		}

		stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
				result.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression. got:%T",
				stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, "x", "==", "y") {
			return
		}

		if len(exp.Consequence.Statements) != 1 {
			t.Fatalf("consequence != 1 statements. got:%d",
				len(result.Statements))
		}

		consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got:%T",
				exp.Consequence.Statements[0])
		}

		if !testLiteralExpression(t, consequence.Expression, 1) {
			return
		}

		if !test.withAlternative && exp.Alternative != nil {
			t.Fatalf("exp.Alternative was not nil. got:%+v", exp.Alternative)
		}

		if !test.withAlternative {
			continue
		}

		// test alternative

		if exp.Alternative == nil {
			t.Fatalf("exp.Alternative was nil. got:%+v", exp.Alternative)
		}

		alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got:%T",
				exp.Alternative.Statements[0])
		}

		if !testLiteralExpression(t, alternative.Expression, 0) {
			return
		}
	}

}

func TestParseIdent(t *testing.T) {
	input := "abcdef"

	l := lexer.New(input)
	p := New(l)
	result := p.ParseProgram()
	checkParserErrors(t, p)

	if len(result.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements. got:%d",
			len(result.Statements))
	}

	stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
			result.Statements[0])
	}

	if !testIdent(t, stmt.Expression, "abcdef") {
		return
	}
}

func TestParseBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		result := p.ParseProgram()
		checkParserErrors(t, p)

		if len(result.Statements) != 1 {
			t.Fatalf("program has incorrect number of statements. got:%d",
				len(result.Statements))
		}

		stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
				result.Statements[0])
		}

		if !testBooleanLiteral(t, stmt.Expression, test.expected) {
			return
		}
	}
}

func TestParseCallExpression(t *testing.T) {
	input := "sum(1, 3 + 12, 4 * 5, a / b)"

	l := lexer.New(input)
	p := New(l)
	result := p.ParseProgram()
	checkParserErrors(t, p)

	if len(result.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements. got:%d",
			len(result.Statements))
	}

	stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
			result.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got:%T",
			result.Statements[0])
	}

	if !testIdent(t, exp.Function, "sum") {
		return
	}

	if len(exp.Arguments) != 4 {
		t.Fatalf("wrong len of arguments. got:%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 3, "+", 12)
	testInfixExpression(t, exp.Arguments[2], 4, "*", 5)
	testInfixExpression(t, exp.Arguments[3], "a", "/", "b")
}

func TestParseFunctionParameters(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		result := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
				result.Statements[0])
		}

		fnExpr, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got:%T",
				stmt.Expression)
		}

		if len(fnExpr.Parameters) != len(test.expectedParams) {
			t.Fatalf("wrong len of parameters. expected: %d got:%d",
				len(test.expectedParams), len(fnExpr.Parameters))
		}

		for i, ident := range test.expectedParams {
			testLiteralExpression(t, fnExpr.Parameters[i], ident)
		}
	}
}

func TestParseFunction(t *testing.T) {
	input := `fn(x, y) { x * y }`

	l := lexer.New(input)
	p := New(l)
	result := p.ParseProgram()
	checkParserErrors(t, p)

	if len(result.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements. got:%d",
			len(result.Statements))
	}

	stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
			result.Statements[0])
	}

	fl, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got:%T",
			result.Statements[0])
	}

	if len(fl.Parameters) != 2 {
		t.Fatalf("wrong len of parameters. expected: 2 got:%d",
			len(fl.Parameters))
	}

	testLiteralExpression(t, fl.Parameters[0], "x")
	testLiteralExpression(t, fl.Parameters[1], "y")

	if len(fl.Body.Statements) != 1 {
		t.Fatalf("fl.Body.Statements has incorrect number of statements. got:%d",
			len(fl.Body.Statements))
	}

	bodyStmt, ok := fl.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("fl.Body.Statements[0] is not ast.ExpressionStatement. got:%T",
			fl.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "*", "y")
}

func TestParseString(t *testing.T) {
	input := `"abc qqqr"`

	l := lexer.New(input)
	p := New(l)
	result := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
			result.Statements[0])
	}

	sl, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.StringLiteral. got:%T",
			result.Statements[0])
	}

	if sl.Value != "abc qqqr" {
		t.Errorf("sl.Value not %q got:%q", "abc qqqr", sl.Value)
	}
}

func TestParseArray(t *testing.T) {
	input := "[5, 25 + 1, a, 5 * 3]"

	l := lexer.New(input)
	p := New(l)
	result := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
			result.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ArrayLiteral. got:%T",
			result.Statements[0])
	}

	if len(array.Elements) != 4 {
		t.Fatalf("wrong len of array not 4. got:%d", len(array.Elements))
	}

	testLiteralExpression(t, array.Elements[0], 5)
	testInfixExpression(t, array.Elements[1], 25, "+", 1)
	testLiteralExpression(t, array.Elements[2], "a")
	testInfixExpression(t, array.Elements[3], 5, "*", 3)
}

func TestParseIndexExpression(t *testing.T) {
	input := "x[5 + 1]"

	l := lexer.New(input)
	p := New(l)
	result := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
			result.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IndexExpression. got:%T",
			result.Statements[0])
	}

	if !testIdent(t, exp.Left, "x") {
		return
	}

	testInfixExpression(t, exp.Index, 5, "+", 1)
}

func TestParseEmptyHashMap(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	result := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
			result.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashMapLiteral. got:%T",
			result.Statements[0])
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("wrong len of hash.Pairs. got:%d expected: 0", len(hash.Pairs))
	}
}

func TestParseHashMapStringKeys(t *testing.T) {
	input := `{"abc": 123, "q": 321}`

	l := lexer.New(input)
	p := New(l)
	result := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
			result.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashMapLiteral. got:%T",
			result.Statements[0])
	}

	expected := map[string]int64{
		"abc": 123,
		"q":   321,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("wrong len of hash.Pairs. got:%d expected: %d", len(hash.Pairs), len(expected))
	}

	for key, value := range hash.Pairs {
		lit, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got:%T", key)
		}

		expectedVal := expected[lit.ToString()]
		testIntLiteral(t, value, expectedVal)
	}
}

func TestParseHashMapBooleanKeys(t *testing.T) {
	input := `{true: 123, false: 321}`

	l := lexer.New(input)
	p := New(l)
	result := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
			result.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashMapLiteral. got:%T",
			result.Statements[0])
	}

	expected := map[string]int64{
		"true":  123,
		"false": 321,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("wrong len of hash.Pairs. got:%d expected: %d", len(hash.Pairs), len(expected))
	}

	for key, value := range hash.Pairs {
		lit, ok := key.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not ast.Boolean. got:%T", key)
		}

		expectedVal := expected[lit.ToString()]
		testIntLiteral(t, value, expectedVal)
	}
}

func TestParseHashMapIntegerKeys(t *testing.T) {
	input := `{1: 123, 2: 321}`

	l := lexer.New(input)
	p := New(l)
	result := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := result.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("result.Statements[0] is not ast.ExpressionStatement. got:%T",
			result.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashMapLiteral. got:%T",
			result.Statements[0])
	}

	expected := map[string]int64{
		"1": 123,
		"2": 321,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("wrong len of hash.Pairs. got:%d expected: %d", len(hash.Pairs), len(expected))
	}

	for key, value := range hash.Pairs {
		lit, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not ast.IntegerLiteral. got:%T", key)
		}

		expectedVal := expected[lit.ToString()]
		testIntLiteral(t, value, expectedVal)
	}
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left any,
	op string,
	right any,
) bool {
	ev, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression not *ast.InfixExpression. got:%T", exp)
		return false
	}

	if !testLiteralExpression(t, ev.Left, left) {
		return false
	}

	if ev.Operator != op {
		t.Errorf("ev.Operator is not %s. got:%s", op, ev.Operator)
		return false
	}

	if !testLiteralExpression(t, ev.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntLiteral(t, exp, int64(v))
	case int64:
		return testIntLiteral(t, exp, v)
	case string:
		return testIdent(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp undefined. got:%T", exp)
	return false
}

func testIntLiteral(t *testing.T, exp ast.Expression, expected int64) bool {
	ev, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expression not *ast.IntegerLiteral. got:%T", exp)
		return false
	}

	if ev.Value != expected {
		t.Errorf("ev.Value is not %d. got:%d", expected, ev.Value)
		return false
	}

	if ev.TokenLiteral() != strconv.FormatInt(expected, 10) {
		t.Errorf("ev.TokenLiteral is not %d. got:%s",
			expected, ev.TokenLiteral())
		return false
	}

	return true
}

func testIdent(t *testing.T, exp ast.Expression, expected string) bool {
	ev, ok := exp.(*ast.Ident)
	if !ok {
		t.Errorf("expression not *ast.Ident. got:%T", exp)
		return false
	}

	if ev.Value != expected {
		t.Errorf("ev.Value is not %s. got:%s", expected, ev.Value)
		return false
	}

	if ev.TokenLiteral() != expected {
		t.Errorf("ev.TokenLiteral is not %s. got:%s",
			expected, ev.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, expected bool) bool {
	ev, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("expression not *ast.Ident. got:%T", exp)
		return false
	}

	if ev.Value != expected {
		t.Errorf("ev.Value is not %t. got:%t", expected, ev.Value)
		return false
	}

	if ev.TokenLiteral() != strconv.FormatBool(expected) {
		t.Errorf("ev.TokenLiteral is not %t. got:%s",
			expected, ev.TokenLiteral())
		return false
	}

	return true
}
