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

		return evalPrefixExpressing(node.Operator, right)

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
		function := Eval(node.FnName, env)
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

func evalPrefixExpressing(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalExclOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s", op)
	}
}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntInfixExpression(op, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(op, left, right)
	case op == "==":
		return boolToBooleanObj(left == right)
	case op == "!=":
		return boolToBooleanObj(left != right)
	case op == "||":
		return boolToBooleanObj(left.(*object.Boolean).Value || right.(*object.Boolean).Value)
	case op == "&&":
		return boolToBooleanObj(left.(*object.Boolean).Value && right.(*object.Boolean).Value)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalExclOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntInfixExpression(op string, left object.Object, right object.Object) object.Object {
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

func evalStringInfixExpression(op string, left object.Object, right object.Object) object.Object {
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

func callFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("call not a function: %s", fn.Type())
	}

	extEnv := extendFuncEnv(function, args)
	ev := Eval(function.Body, extEnv)
	return unwrapReturn(ev)
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
