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

func Eval(n ast.Node) object.Object {
	switch node := n.(type) {
	case *ast.Program:
		return evalProgram(node)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.ReturnStatement:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}

		return &object.Return{Value: val}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		if node.Value {
			return TRUE
		}

		return FALSE

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalPrefixExpressing(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		switch r := result.(type) {
		case *object.Return:
			return r.Value
		case *object.Error:
			return r
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)
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

func evalIfExpression(node *ast.IfExpression) object.Object {
	condition := Eval(node.Condition)
	if isError(condition) {
		return condition
	}

	if condition != TRUE && condition != FALSE {
		return newError("non boolean condition in if statement")
	}

	if condition == TRUE {
		return Eval(node.Consequence)
	} else if node.Alternative != nil {
		return Eval(node.Alternative)
	} else {
		return NULL
	}
}
