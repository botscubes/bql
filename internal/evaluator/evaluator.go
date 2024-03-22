package evaluator

import (
	"fmt"

	"github.com/botscubes/bql/internal/ast"
	"github.com/botscubes/bql/internal/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func newError(formating string, parameters ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(formating, parameters...)}
}

func Eval(n ast.Node) object.Object {
	switch node := n.(type) {
	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		if node.Value {
			return TRUE
		}

		return FALSE

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpressing(node.Operator, right)
	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)
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
