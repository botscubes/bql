package evaluator

import (
	"fmt"

	"github.com/botscubes/bql/internal/ast"
	"github.com/botscubes/bql/internal/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func newError(formating string, parameters ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(formating, parameters...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func boolToBooleanObj(b bool) *object.Boolean {
	if b {
		return TRUE
	}

	return FALSE
}

func Eval(n ast.Node, env *object.Env) object.Object {
	switch node := n.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		return &object.Return{Value: val}

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		if node.Value {
			return TRUE
		}

		return FALSE

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Ident:
		return evalIdent(node, env)

	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return callFunction(function, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.HashMapLiteral:
		return evalHashMap(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Env) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch r := result.(type) {
		case *object.Return:
			return r.Value
		case *object.Error:
			return r
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Env) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalExclOpExpr(right)
	case "-":
		return evalMinusPrefixOpExpr(right)
	default:
		return newError("unknown operator: %s", op)
	}
}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntInfixExpr(op, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpr(op, left, right)
	case op == "==":
		return boolToBooleanObj(left == right)
	case op == "!=":
		return boolToBooleanObj(left != right)
	case op == "||":
		return boolToBooleanObj(left.(*object.Boolean).Value || right.(*object.Boolean).Value)
	case op == "&&":
		return boolToBooleanObj(left.(*object.Boolean).Value && right.(*object.Boolean).Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalExclOpExpr(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOpExpr(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntInfixExpr(op string, left object.Object, right object.Object) object.Object {
	lVal := left.(*object.Integer).Value
	rVal := right.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: lVal + rVal}
	case "-":
		return &object.Integer{Value: lVal - rVal}
	case "*":
		return &object.Integer{Value: lVal * rVal}
	case "/":
		return &object.Integer{Value: lVal / rVal}
	case "%":
		return &object.Integer{Value: lVal % rVal}
	case "==":
		return boolToBooleanObj(lVal == rVal)
	case "!=":
		return boolToBooleanObj(lVal != rVal)
	case "<":
		return boolToBooleanObj(lVal < rVal)
	case ">":
		return boolToBooleanObj(lVal > rVal)
	case "<=":
		return boolToBooleanObj(lVal <= rVal)
	case ">=":
		return boolToBooleanObj(lVal >= rVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalStringInfixExpr(op string, left object.Object, right object.Object) object.Object {
	lVal := left.(*object.String).Value
	rVal := right.(*object.String).Value
	switch op {
	case "+":
		return &object.String{Value: lVal + rVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIfExpression(node *ast.IfExpression, env *object.Env) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if condition != TRUE && condition != FALSE {
		return newError("non boolean condition in if statement")
	}

	if condition == TRUE {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdent(node *ast.Ident, env *object.Env) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exprs []ast.Expression, env *object.Env) []object.Object {
	var result []object.Object

	for _, exp := range exprs {
		ev := Eval(exp, env)
		if isError(ev) {
			return []object.Object{ev}
		}
		result = append(result, ev)
	}
	return result
}

func callFunction(function object.Object, args []object.Object) object.Object {
	switch fn := function.(type) {
	case *object.Function:
		extEnv := extendFuncEnv(fn, args)
		ev := Eval(fn.Body, extEnv)
		return unwrapReturn(ev)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("call not a function: %s", fn.Type())
	}
}

func extendFuncEnv(fn *object.Function, args []object.Object) *object.Env {
	env := object.NewEnclosedEnv(fn.Env)

	for id, param := range fn.Parameters {
		env.Set(param.Value, args[id])
	}

	return env
}

func unwrapReturn(obj object.Object) object.Object {
	if val, ok := obj.(*object.Return); ok {
		return val.Value
	}

	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExp(left, index)
	case left.Type() == object.HASH_MAP_OBJ:
		return evalHashMapIndexExp(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExp(left, index object.Object) object.Object {
	array := left.(*object.Array)
	idx := index.(*object.Integer).Value

	if idx < 0 || idx > int64(len(array.Elements)-1) {
		return NULL
	}

	return array.Elements[idx]
}

func evalHashMap(node *ast.HashMapLiteral, env *object.Env) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for knode, vnode := range node.Pairs {
		key := Eval(knode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(vnode, env)
		if isError(value) {
			return value
		}

		pairs[hashKey.HashKey()] = object.HashPair{Key: key, Value: value}
	}

	return &object.HashMap{Pairs: pairs}
}

func evalHashMapIndexExp(hashMap, index object.Object) object.Object {
	hashObject := hashMap.(*object.HashMap)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}
