package evaluator

import (
	"github.com/botscubes/bql/internal/ast"
	"github.com/botscubes/bql/internal/object"
)

func Eval(n ast.Node) object.Object {
	switch node := n.(type) {
	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
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
